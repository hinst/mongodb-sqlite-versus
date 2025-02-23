package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"
)

const DB_FILE_PATH = "./test-sqlite.db"

type SqliteTest struct {
}

func (me *SqliteTest) initialize() {
	if checkFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var setupText = readStringFromFile(executablePath + "/setup.sql")
	assertResultError(db.Exec(setupText))
	db.Close()
}

func (me *SqliteTest) testInsertion(users []*User, threadCount int) time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for i := 0; i < threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			writeSqlite(usersChannel)
		}()
	}
	for _, user := range users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)
	return elapsed
}

func (me *SqliteTest) testReading(users []*User, threadCount int) time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for i := 0; i < threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			readSqlite(usersChannel)
		}()
	}
	for _, user := range users {
		usersChannel <- user
	}
	var elapsed = time.Since(beginning)
	return elapsed
}

func (me *SqliteTest) compress() (int64, int64) {
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var sizeBeforeVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	db.Exec("VACUUM;")
	db.Close()

	var sizeAfterVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	return sizeBeforeVacuum, sizeAfterVacuum
}

func (me *SqliteTest) run(users []*User, threadCount int) {
	me.initialize()

	var insertDuration = me.testInsertion(users, threadCount)
	var insertionsPerSecond = float64(len(users)) / insertDuration.Seconds()

	var readDuration = me.testReading(users, threadCount)
	var readsPerSecond = float64(len(users)) / readDuration.Seconds()

	var sizeBeforeVacuum, sizeAfterVacuum = me.compress()
	fmt.Printf("SQLite file size: %v -> %v\n",
		formatFileSize(sizeBeforeVacuum), formatFileSize(sizeAfterVacuum))
	fmt.Printf(TAB+"insertion duration: %v, rows per second: %.1f\n", insertDuration, insertionsPerSecond)
	fmt.Printf(TAB+"reading duration: %v, rows per second: %.1f\n", readDuration, readsPerSecond)
}

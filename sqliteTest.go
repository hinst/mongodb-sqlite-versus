package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

const DB_FILE_PATH = "./test-sqlite.db"

type SqliteTest struct {
	users       []*User
	batchSize   int
	threadCount int
}

func (me *SqliteTest) prepare() {
	if checkFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var setupText = readStringFromFile(executablePath + "/setup.sql")
	assertResultError(db.Exec(setupText))
	db.Close()
}

func (me *SqliteTest) run() {
	me.prepare()

	var insertDuration = me.testInsertion()
	var insertionsPerSecond = float64(len(me.users)) / insertDuration.Seconds()

	var readDuration = me.testReading()
	var readsPerSecond = float64(len(me.users)) / readDuration.Seconds()

	var sizeBeforeVacuum, sizeAfterVacuum = me.compress()
	fmt.Printf("SQLite file size: %v -> %v\n",
		formatFileSize(sizeBeforeVacuum), formatFileSize(sizeAfterVacuum))
	fmt.Printf(TAB+"insertion duration: %v, rows per second: %v\n",
		insertDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))
	fmt.Printf(TAB+"reading duration: %v, rows per second: %v\n",
		readDuration, humanize.CommafWithDigits(readsPerSecond, 0))
}

func (me *SqliteTest) testInsertion() time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for i := 0; i < me.threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			writeSqlite(usersChannel, me.batchSize)
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)
	return elapsed
}

func (me *SqliteTest) testReading() time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for i := 0; i < me.threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			readSqlite(usersChannel, me.batchSize)
		}()
	}
	for _, user := range me.users {
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

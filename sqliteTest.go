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
	assertError(db.Close())
}

func (me *SqliteTest) run() {
	me.prepare()

	var insertDuration = me.runInserts()
	var insertionsPerSecond = float64(len(me.users)) / insertDuration.Seconds()

	var readDuration = me.runQueries()
	var readsPerSecond = float64(len(me.users)) / readDuration.Seconds()

	var beginning = time.Now()
	var sizeBefore, sizeAfter = me.compress()
	var compressionDuration = time.Since(beginning)

	fmt.Printf("SQLite file size: %v -> %v, compression duration: %v\n",
		formatFileSize(sizeBefore), formatFileSize(sizeAfter), compressionDuration)
	fmt.Printf(TAB+"insertion duration: %v, rows per second: %v\n",
		insertDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))
	fmt.Printf(TAB+"reading duration: %v, rows per second: %v\n",
		readDuration, humanize.CommafWithDigits(readsPerSecond, 0))
}

func (me *SqliteTest) runInserts() time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for range me.threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.writeUsers(usersChannel)
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

func (me *SqliteTest) runQueries() time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for i := 0; i < me.threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.readUsers(usersChannel)
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	var elapsed = time.Since(beginning)
	return elapsed
}

func (me *SqliteTest) readUsers(users chan *User) {
	var db *sql.DB
	var counter = 0
	for user := range users {
		if nil == db {
			db = assertResultError(sql.Open("sqlite3", "file:"+DB_FILE_PATH+"?mode=ro"))
		}
		var row = db.QueryRow("SELECT name, passwordHash, accessToken, email, createdAt, level FROM users WHERE id=?", user.SqliteId)
		assertError(row.Err())
		var userB User = User{SqliteId: user.SqliteId}
		var createdAt int64
		row.Scan(&userB.Name, &userB.PasswordHash, &userB.AccessToken, &userB.Email, &createdAt, &userB.Level)
		userB.CreatedAt = time.Unix(createdAt, 0)
		// if *user != userB {
		// 	fmt.Printf("Expected [%v] but got [%v]\n", *user, userB)
		// }
		assertCondition(*user == userB, "Users must be equal")
		counter += 1
		if (counter%me.batchSize) == 0 && db != nil {
			assertError(db.Close())
			db = nil
		}
	}
	if db != nil {
		assertError(db.Close())
	}
}

func (me *SqliteTest) writeUsers(users chan *User) {
	var db *sql.DB
	var counter = 0
	for user := range users {
		if nil == db {
			db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
		}
		var row = db.QueryRow("INSERT INTO users (name, passwordHash, accessToken, email, createdAt, level) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
			user.Name, user.PasswordHash, user.AccessToken, user.Email, user.CreatedAt.Unix(), user.Level)
		assertError(row.Scan(&user.SqliteId))
		counter += 1
		if (counter%me.batchSize) == 0 && db != nil {
			assertError(db.Close())
			db = nil
		}
	}
	if db != nil {
		assertError(db.Close())
	}
}

func (me *SqliteTest) compress() (int64, int64) {
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var sizeBeforeVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	assertResultError(db.Exec("VACUUM;"))
	assertError(db.Close())

	var sizeAfterVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	return sizeBeforeVacuum, sizeAfterVacuum
}

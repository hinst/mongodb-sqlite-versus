package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const DB_FILE_PATH = "./test-sqlite.db"

func testSqlite(users []*User, threadCount int) {
	if checkFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var setupText = readStringFromFile(executablePath + "/setup.sql")
	assertResultError(db.Exec(setupText))
	db.Close()

	var beginning = time.Now()
	var usersChannel = make(chan *User)
	for i := 0; i < threadCount; i++ {
		go writeSqlite(usersChannel)
	}
	for _, user := range users {
		usersChannel <- user
	}
	close(usersChannel)
	var elapsed = time.Since(beginning)

	db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var sizeBeforeVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	db.Exec("VACUUM;")
	db.Close()

	var sizeAfterVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	fmt.Printf("SQLite time: %v, file size: %v -> %v\n", elapsed,
		formatFileSize(sizeBeforeVacuum), formatFileSize(sizeAfterVacuum))
}

package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const DB_FILE_PATH = "./test-sqlite.db"

func testSqlite(rowCount int) {
	if CheckFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}

	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	defer db.Close()
	var setupText = ReadStringFromFile("./setup.sql")
	assertResultError(db.Exec(setupText))
	var users = generateRandomUsers(rowCount)

	var beginning = time.Now()
	for i := 0; i < rowCount; i++ {
		assertResultError(db.Exec("INSERT INTO users (name, passwordHash, email, createdAt, level) VALUES (?, ?, ?, ?, ?)",
			users[i].name, users[i].passwordHash, users[i].email, users[i].createdAt, users[i].level))
	}
	var elapsed = time.Since(beginning)
	fmt.Printf("Inserted. Rows: [%d], time: %v\n", rowCount, elapsed)

	db.Exec("VACUUM;")
}

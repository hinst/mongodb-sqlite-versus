package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func testSqlite(users []*User) {
	const DB_FILE_PATH = "./test-sqlite.db"
	if checkFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var setupText = readStringFromFile(executablePath + "/setup.sql")
	assertResultError(db.Exec(setupText))

	var beginning = time.Now()
	for i := 0; i < len(users); i++ {
		assertResultError(db.Exec("INSERT INTO users (name, passwordHash, email, createdAt, level) VALUES (?, ?, ?, ?, ?)",
			users[i].Name, users[i].PasswordHash, users[i].Email, users[i].CreatedAt, users[i].Level))
	}
	var elapsed = time.Since(beginning)
	db.Close()

	db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	var sizeBeforeVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	db.Exec("VACUUM;")
	db.Close()

	var sizeAfterVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
	fmt.Printf("SQLite time: %v, file size: %v -> %v\n", elapsed,
		formatFileSize(sizeBeforeVacuum), formatFileSize(sizeAfterVacuum))
}

package main

import (
	"database/sql"
)

func writeSqlite(users chan *User, batchSize int) {
	var db *sql.DB
	defer func() {
		if db != nil {
			db.Close()
			db = nil
		}
	}()
	var counter = 0
	for {
		if nil == db {
			db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
		}
		var user, ok = <-users
		if !ok {
			break
		}
		var row = db.QueryRow("INSERT INTO users (name, passwordHash, accessToken, email, createdAt, level) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
			user.Name, user.PasswordHash, user.AccessToken, user.Email, user.CreatedAt.Unix(), user.Level)
		assertError(row.Scan(&user.SqliteId))
		counter += 1
		if (counter%batchSize) == 0 && db != nil {
			db.Close()
			db = nil
		}
	}
}

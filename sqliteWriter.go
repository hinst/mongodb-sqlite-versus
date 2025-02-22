package main

import (
	"database/sql"
)

func writeSqlite(users chan *User) {
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	defer db.Close()
	for {
		var user, ok = <-users
		if !ok {
			break
		}
		assertResultError(db.Exec("INSERT INTO users (name, passwordHash, email, createdAt, level) VALUES (?, ?, ?, ?, ?)",
			user.Name, user.PasswordHash, user.Email, user.CreatedAt, user.Level))
	}
}

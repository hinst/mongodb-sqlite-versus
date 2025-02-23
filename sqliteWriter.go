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
		var row = db.QueryRow("INSERT INTO users (name, passwordHash, accessToken, email, createdAt, level) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
			user.Name, user.PasswordHash, user.AccessToken, user.Email, user.CreatedAt, user.Level)
		assertError(row.Scan(&user.SqliteId))
	}
}

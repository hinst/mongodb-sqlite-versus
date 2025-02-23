package main

import (
	"database/sql"
	"time"
)

func readSqlite(users chan *User, batchSize int) {
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
		if (counter%batchSize) == 0 && db != nil {
			db.Close()
			db = nil
		}
	}
}

package main

import "database/sql"

func readSqlite(users []*User) {
	var db = assertResultError(sql.Open("sqlite3", DB_FILE_PATH))
	defer db.Close()
	for _, user := range users {
		var row = db.QueryRow("SELECT (name, passwordHash, accessToken, email, createdAt, level) FROM users WHERE id=?", user)
		var userB User = User{SqliteId: user.SqliteId}
		row.Scan(&userB.Name, &userB.PasswordHash, &userB.AccessToken, &userB.Email, &userB.CreatedAt, &userB.Level)
		assertCondition(*user == userB, "Users are equal")
	}
}

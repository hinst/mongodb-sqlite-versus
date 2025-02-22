package main

import (
	"database/sql"
	"os"
)

func testSqlite() {
	var db = assertResultError(sql.Open("sqlite3", "./test-sqlite.db"))
	defer db.Close()
	var setupText = string(assertResultError(os.ReadFile("./setup.sql")))
	assertResultError(db.Exec(setupText))
}

package main

import (
	"database/sql"
	_ "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func test() {
	var db = AssertResultError(sql.Open("sqlite3", "./test-sqlite.db"))
}

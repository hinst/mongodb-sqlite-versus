package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const rowCount = 1_000_000
	fmt.Println("STARTING")
	testSqlite(rowCount)
}

package main

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const rowCount = 100_000
	var users = generateRandomUsers(rowCount)
	fmt.Printf("Testing rows [%v]\n", len(users))

	if false {
		time.Sleep(100 * time.Millisecond)
		testJson(users)
	}

	if false {
		time.Sleep(100 * time.Millisecond)
		testSqlite(users)
	}

	if true {
		time.Sleep(100 * time.Millisecond)
		testMongo(users)
	}
}

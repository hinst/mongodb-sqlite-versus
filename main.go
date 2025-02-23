package main

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const rowCount = 1
	const threadCount = 4
	var users = generateRandomUsers(rowCount)
	fmt.Printf("Testing rows [%v]\n", len(users))

	time.Sleep(100 * time.Millisecond)
	testJson(users)

	time.Sleep(100 * time.Millisecond)
	new(SqliteTest).run(users, threadCount)

	time.Sleep(100 * time.Millisecond)
	testMongo(users, threadCount)
}

package main

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const rowCount = 100_000
	const batchSize = 10
	const threadCount = 4
	const waitInterval = 100 * time.Millisecond
	var users = generateRandomUsers(rowCount)
	fmt.Printf("Testing rows[%v] threads[%v]\n", len(users), threadCount)

	time.Sleep(waitInterval)
	testJson(users)

	time.Sleep(waitInterval)
	(&SqliteTest{users, batchSize, threadCount}).run()

	time.Sleep(waitInterval)
	(&MongoTest{users, batchSize, threadCount}).run()
}

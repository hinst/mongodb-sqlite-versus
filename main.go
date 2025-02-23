package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const rowCount = 100_000
	const batchSize = 100
	const threadCount = 4
	const waitInterval = 1 * time.Second
	var users = generateRandomUsers(rowCount)
	fmt.Printf("Testing rows[%v] batchSize[%v] threads[%v]\n",
		humanize.Comma(int64(len(users))), batchSize, threadCount)

	time.Sleep(waitInterval)
	testJson(users)

	time.Sleep(waitInterval)
	(&SqliteTest{users, batchSize, threadCount}).run()

	time.Sleep(waitInterval)
	(&MongoTest{users, batchSize, threadCount}).run()
}

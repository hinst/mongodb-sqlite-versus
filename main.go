package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
)

var testers = map[string](func(users []*User, batchSize int, threadCount int)){
	"json": func(users []*User, batchSize int, threadCount int) {
		testJson(users)
	},
	"sqlite": func(users []*User, batchSize int, threadCount int) {
		(&SqliteTest{users, batchSize, threadCount}).run()
	},
	"mongo": func(users []*User, batchSize int, threadCount int) {
		(&MongoTest{users, batchSize, threadCount}).run()
	},
}

func getTesterKeys() string {
	var keys = make([]string, 0, len(testers))
	for key := range testers {
		keys = append(keys, key)
	}
	return fmt.Sprintf("%v", keys)
}

func main() {
	const batchSize = 100
	const threadCount = 4
	var testKey = flag.String("test", "sqlite", getTesterKeys())
	var rowCount = flag.Int("rows", 1_000_000, "number of rows to test")
	flag.Parse()
	var tester = testers[*testKey]
	assertCondition(tester != nil, "Invalid tester key: "+*testKey)

	var users = generateRandomUsers(*rowCount)
	fmt.Printf("Testing rows[%v] batchSize[%v] threads[%v] db[%v]\n",
		humanize.Comma(int64(len(users))), batchSize, threadCount, *testKey)
	var beginning = time.Now()
	tester(users, batchSize, threadCount)
	var elapsed = time.Since(beginning)
	fmt.Printf("Test complete %v\n", elapsed)
}

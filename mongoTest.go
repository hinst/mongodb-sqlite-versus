package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_DB_URL = "mongodb://localhost:27017"

type MongoTest struct {
	users       []*User
	batchSize   int
	threadCount int
}

func (me *MongoTest) prepare() {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	defer func() { assertError(client.Disconnect(context.Background())) }()
	client.Database("test").Drop(context.Background())
}

func (me *MongoTest) run() {
	me.prepare()
	var insertionsDuration = me.runInsertions()
	var insertionsPerSecond = float64(len(me.users)) / insertionsDuration.Seconds()

	var beginning = time.Now()
	var sizeBefore, sizeAfter = me.compress()
	var compressionDuration = time.Since(beginning)

	fmt.Printf("MongoDB file size: %v -> %v, compression duration %v\n",
		formatFileSize(sizeBefore), formatFileSize(sizeAfter), compressionDuration)
	fmt.Printf(TAB+"insertion duration: %v, rows per second: %v\n",
		insertionsDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))
}

func (me *MongoTest) runInsertions() time.Duration {
	var beginning = time.Now()
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup
	for i := 0; i < me.threadCount; i++ {
		waitGroup.Add(1)
		go func() {
			writeUsers(usersChannel, me.batchSize)
			waitGroup.Done()
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)
	return elapsed
}

func (me *MongoTest) compress() (sizeBefore int64, sizeAfter int64) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	defer func() { assertError(client.Disconnect(context.Background())) }()
	var db = client.Database("test")
	db.RunCommand(context.Background(), bson.M{"compact": "users"})
	var stats = getMongoDbStats(db)
	return int64(stats.dataSize), int64(stats.storageSize)
}

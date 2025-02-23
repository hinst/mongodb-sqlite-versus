package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var queriesDuration = me.runQueries()
	var queriesPerSecond = float64(len(me.users)) / queriesDuration.Seconds()

	var beginning = time.Now()
	var sizeBefore, sizeAfter = me.compress()
	var compressionDuration = time.Since(beginning)

	fmt.Printf("MongoDB file size: %v -> %v, compression duration %v\n",
		formatFileSize(sizeBefore), formatFileSize(sizeAfter), compressionDuration)
	fmt.Printf(TAB+"insertions duration: %v, rows per second: %v\n",
		insertionsDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))
	fmt.Printf(TAB+"queries duration: %v, rows per second: %v\n",
		queriesDuration, humanize.CommafWithDigits(queriesPerSecond, 0))
}

func (me *MongoTest) runInsertions() time.Duration {
	var beginning = time.Now()
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup
	for range me.threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.writeUsers(usersChannel)
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

func (me *MongoTest) writeUsers(users chan *User) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client *mongo.Client
	var counter = 0
	defer func() {
		if client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}()
	for user := range users {
		if nil == client {
			client = assertResultError(mongo.Connect(context.Background(), clientOptions))
		}
		var db = client.Database("test")
		var result = assertResultError(db.Collection("users").InsertOne(context.Background(), user))
		user._id = result.InsertedID.(primitive.ObjectID)
		counter += 1
		if (counter%me.batchSize) == 0 && client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}
}

func (me *MongoTest) runQueries() time.Duration {
	var beginning = time.Now()
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup
	for range me.threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.readUsers(usersChannel, me.batchSize)
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

func (me *MongoTest) readUsers(usersChannel chan *User, batchSize int) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client *mongo.Client
	defer func() {
		if client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}()
	var counter = 0
	for user := range usersChannel {
		if nil == client {
			client = assertResultError(mongo.Connect(context.Background(), clientOptions))
		}
		var collection = client.Database("test").Collection("users")
		var result = collection.FindOne(context.Background(), bson.M{"_id": user._id})
		assertError(result.Err())
		counter += 1
		if (counter%batchSize) == 0 && client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}
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

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

const MONGO_DB_URL = "mongodb://localhost:27017/?keepAlive=false"

type MongoTest struct {
	users       []*User
	batchSize   int
	threadCount int
}

func (me *MongoTest) prepare() {
	var client = me.open()
	defer me.close(client)
	client.Database("test").Drop(context.Background())
}

func (me *MongoTest) open() *mongo.Client {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	return assertResultError(mongo.Connect(context.Background(), clientOptions))
}

func (me *MongoTest) run() {
	me.prepare()
	var insertionsDuration = me.runInsertions()
	var insertionsPerSecond = float64(len(me.users)) / insertionsDuration.Seconds()
	fmt.Printf(TAB+"insertions duration: %v, rows per second: %v\n",
		insertionsDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))

	var queriesDuration = me.runQueries(me.threadCount)
	var queriesPerSecond = float64(len(me.users)) / queriesDuration.Seconds()
	fmt.Printf(TAB+"queries duration: %v, rows per second: %v\n",
		queriesDuration, humanize.CommafWithDigits(queriesPerSecond, 0))

	var combinedReadDuration, combinedUpdateDuration = me.runCombined()
	var combinedReadsPerSecond = float64(len(me.users)) / combinedReadDuration.Seconds()
	var combinedUpdatesPerSecond = float64(len(me.users)) / combinedUpdateDuration.Seconds()

	fmt.Printf(TAB+"combined read & update benchmark: %v reads per second, %v updates per second\n",
		humanize.CommafWithDigits(combinedReadsPerSecond, 0),
		humanize.CommafWithDigits(combinedUpdatesPerSecond, 0))
	fmt.Printf(TAB+TAB+"read duration %v, update duration %v\n",
		combinedReadDuration, combinedUpdateDuration)

	var beginning = time.Now()
	var sizeBefore, sizeAfter = me.compress()
	var compressionDuration = time.Since(beginning)

	fmt.Printf("MongoDB file size: %v -> %v, compression duration %v\n",
		formatFileSize(sizeBefore), formatFileSize(sizeAfter), compressionDuration)
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
	var client *mongo.Client
	var counter = 0
	defer func() { me.close(client) }()
	for user := range users {
		if nil == client {
			client = me.open()
		}
		var db = client.Database("test")
		var result = assertResultError(db.Collection("users").InsertOne(context.Background(), user))
		user.MongoId = result.InsertedID.(primitive.ObjectID)
		counter += 1
		if (counter%me.batchSize) == 0 && client != nil {
			client = me.close(client)
		}
	}
}

func (me *MongoTest) runQueries(threadCount int) time.Duration {
	var beginning = time.Now()
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup
	for range threadCount {
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
	var client *mongo.Client
	defer func() { me.close(client) }()
	var counter = 0
	for user := range usersChannel {
		if nil == client {
			client = me.open()
		}
		var collection = client.Database("test").Collection("users")
		var result = collection.FindOne(context.Background(), bson.M{"_id": user.MongoId})
		assertError(result.Err())
		var userB User
		result.Decode(&userB)
		assertCondition(user.compare(&userB), "Users must be equal")
		counter += 1
		if (counter%batchSize) == 0 && client != nil {
			client = me.close(client)
		}
	}
}

func (me *MongoTest) runCombined() (readDuration time.Duration, updateDuration time.Duration) {
	var threadCount = max(1, me.threadCount/2)
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		updateDuration = me.runUpdates(threadCount)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		readDuration = me.runQueries(threadCount)
	}()

	waitGroup.Wait()
	return readDuration, updateDuration
}

func (me *MongoTest) runUpdates(threadCount int) time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for range threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.updateUsers(usersChannel)
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

func (me *MongoTest) updateUsers(users chan *User) {
	var client *mongo.Client
	defer func() { me.close(client) }()
	var counter = 0
	for user := range users {
		if nil == client {
			client = me.open()
		}
		var collection = client.Database("test").Collection("users")
		assertResultError(collection.UpdateOne(
			context.Background(),
			bson.M{"_id": user.MongoId},
			bson.M{"$set": user},
		))
		counter += 1
		if (counter%me.batchSize) == 0 && client != nil {
			client = me.close(client)
		}
	}
}

func (me *MongoTest) compress() (sizeBefore int64, sizeAfter int64) {
	var client = me.open()
	defer me.close(client)
	var db = client.Database("test")
	db.RunCommand(context.Background(), bson.M{"compact": "users"})
	var stats = getMongoDbStats(db)
	return int64(stats.dataSize), int64(stats.storageSize)
}

func (me *MongoTest) close(client *mongo.Client) *mongo.Client {
	if client != nil {
		assertError(client.Disconnect(context.Background()))
	}
	return nil
}

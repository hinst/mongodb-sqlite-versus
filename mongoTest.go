package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func testMongo(users []*User, threadCount int) {
	var beginning = time.Now()
	var usersChannel = make(chan *User)
	for i := 0; i < threadCount; i++ {
		go func() {
			writeUsers(usersChannel)
		}()
	}
	for _, user := range users {
		usersChannel <- user
	}
	close(usersChannel)
	var elapsed = time.Since(beginning)

	var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	var client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	defer func() {
		assertError(client.Disconnect(context.Background()))
	}()
	var db = client.Database("test")
	db.RunCommand(context.Background(), bson.M{"compact": "users"})
	time.Sleep(10 * time.Second)
	var stats = getMongoDbStats(db)

	fmt.Printf("MongoDB time: %v, size: %v -> %v\n", elapsed,
		formatFileSize(int64(stats.dataSize)), formatFileSize(int64(stats.storageSize)))
}

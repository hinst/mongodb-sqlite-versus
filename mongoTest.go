package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func testMongo(users []*User) {
	var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	client.Database("test").Drop(context.Background())
	time.Sleep(100 * time.Millisecond)

	var beginning = time.Now()
	var db = client.Database("test")
	for _, user := range users {
		assertResultError(db.Collection("users").InsertOne(context.Background(), user))
	}
	assertError(client.Disconnect(context.TODO()))
	var elapsed = time.Since(beginning)

	client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	db = client.Database("test")
	db.RunCommand(context.Background(), bson.M{"compact": "users"})
	time.Sleep(10 * time.Second)
	var stats = getMongoDbStats(db)
	client.Disconnect(context.Background())

	fmt.Printf("MongoDB time: %v, size: %v -> %v\n", elapsed,
		formatFileSize(int64(stats.dataSize)), formatFileSize(int64(stats.storageSize)))
}

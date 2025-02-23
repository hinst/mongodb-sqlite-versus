package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func writeUsers(users chan *User, batchSize int) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client *mongo.Client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	defer func() {
		if client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}()
	for {
		var user, ok = <-users
		if !ok {
			break
		}
		var db = client.Database("test")
		assertResultError(db.Collection("users").InsertOne(context.Background(), user))
	}
}

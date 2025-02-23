package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func writeUsers(users chan *User, batchSize int) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	var client *mongo.Client
	var counter = 0
	defer func() {
		if client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}()
	for {
		if nil == client {
			client = assertResultError(mongo.Connect(context.Background(), clientOptions))
		}
		var user, ok = <-users
		if !ok {
			break
		}
		var db = client.Database("test")
		assertResultError(db.Collection("users").InsertOne(context.Background(), user))
		counter += 1
		if (counter%batchSize) == 0 && client != nil {
			assertError(client.Disconnect(context.Background()))
			client = nil
		}
	}
}

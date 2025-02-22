package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func writeUsers(users chan *User) {
	var clientOptions = options.Client().ApplyURI(MONGO_DB_URL)
	clientOptions.Compressors = []string{"snappy"}
	var client = assertResultError(mongo.Connect(context.Background(), clientOptions))
	defer func() {
		assertError(client.Disconnect(context.Background()))
	}()
	client.Database("test").Drop(context.Background())
	var db = client.Database("test")
	for {
		var user, ok = <-users
		if !ok {
			break
		}
		assertResultError(db.Collection("users").InsertOne(context.Background(), user))
	}
}

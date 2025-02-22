package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getMongoDbSize(db *mongo.Database) int64 {
	var statsResult = db.RunCommand(context.Background(), bson.M{"dbStats": 1})
	var stats bson.M
	assertError(statsResult.Decode(&stats))
	var storageSize = stats["storageSize"].(float64)
	return int64(storageSize)
}

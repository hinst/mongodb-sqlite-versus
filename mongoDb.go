package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbStats struct {
	dataSize    float64
	storageSize float64
}

func getMongoDbStats(db *mongo.Database) MongoDbStats {
	var statsResult = db.RunCommand(context.Background(), bson.M{"dbStats": 1})
	var stats bson.M
	assertError(statsResult.Decode(&stats))
	return MongoDbStats{
		dataSize:    stats["dataSize"].(float64),
		storageSize: stats["storageSize"].(float64),
	}
}

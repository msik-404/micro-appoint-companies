package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getURI() string {
	return fmt.Sprintf("mongodb://%s:%s@mongodb:27017", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

func ConnectDB() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(getURI()).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create a new client and connect to the server
	return mongo.Connect(ctx, opts)
}

func CreateDBIndexes(db *mongo.Database) ([]string, error) {
	coll := db.Collection("companies")
	index := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"type": 1},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return coll.Indexes().CreateMany(ctx, index)
}

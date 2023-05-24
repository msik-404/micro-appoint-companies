package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service *Service) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("services")
	return GenericInsertOne(coll, service)
}

func (service *Service) UpdateOne(db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection("services")
	filter := bson.D{{"_id", service.ID}}
	return GenericUpdateOne(coll, filter, service)
}

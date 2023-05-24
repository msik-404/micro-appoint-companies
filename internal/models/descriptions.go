package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (description *Description) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("descriptions")
	return GenericInsertOne(coll, description)
}

func (description *Description) UpdateOne(db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection("descriptions")
	filter := bson.D{{"_id", description.CompanyID}}
	return GenericUpdateOne(coll, filter, description)
}

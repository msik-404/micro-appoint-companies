package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoModel interface {
	InsertOne(*mongo.Database) (*mongo.InsertOneResult, error)
	UpdateOne(*mongo.Database) (*mongo.UpdateResult, error)
}

type Service struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Price       uint               `json:"price" bson:"price,omitempty"`
	Duration    uint               `json:"duration" bson:"duration,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id,omitempty"`
}

type Description struct {
	CompanyID   primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
}

type Company struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name,omitempty"`
	Type             string             `json:"type" bson:"type,omitempty"`
	Localisation     string             `json:"localisation" bson:"localisation,omitempty"`
	ShortDescription string             `json:"short_description" bson:"short_description,omitempty"`
}

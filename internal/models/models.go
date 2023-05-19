package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Name        string             `json:"name" bson:"name"`
	Price       uint               `json:"price" bson:"price"`
	Duration    uint               `json:"duration" bson:"duration"`
	Description string             `json:"description" bson:"description"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
}

type Description struct {
	CompanyID   primitive.ObjectID `json:"-" bson:"_id"`
	Description string             `json:"description" bson:"description"`
}

type Company struct {
	// should be unique and indexed
	Name string `json:"name" bson:"name"`
	// should be indexed and from some pool of options
	Type             string `json:"type" bson:"type"`
	Localisation     string `json:"localisation" bson:"localisation"`
	ShortDescription string `json:"short_description" bson:"short_description"`
}

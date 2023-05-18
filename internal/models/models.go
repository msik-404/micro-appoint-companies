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

type Company struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

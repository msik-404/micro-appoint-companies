package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Name             string `json:"name" bson:"name"`
	Type             string `json:"type" bson:"type"`
	Localisation     string `json:"localisation" bson:"localisation"`
	ShortDescription string `json:"short_description" bson:"short_description"`
}

// structure which will be received by PostCompaniesEndPoint
type CompanyPost struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Localisation     string `json:"localisation"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
}

func (companyPost *CompanyPost) Insert(db *mongo.Database) (*mongo.InsertOneResult, error) {
	company := Company{
		companyPost.Name,
		companyPost.Type,
		companyPost.Localisation,
		companyPost.ShortDescription,
	}
	coll := db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.InsertOne(ctx, company)
	if err != nil {
		return result, err
	}
	description := Description{
		result.InsertedID.(primitive.ObjectID),
		companyPost.Description,
	}
	coll = db.Collection("descriptions")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.InsertOne(ctx, description)
}

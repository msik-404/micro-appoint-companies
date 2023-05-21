package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoModel interface {
	InsertOne(*mongo.Database) (*mongo.InsertOneResult, error)
}

func genericInsertOne[T MongoModel](coll *mongo.Collection, item T) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.InsertOne(ctx, item)
}

func GenericFindOne[T any](coll *mongo.Collection, filter bson.D) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var doc T
	err := coll.FindOne(ctx, filter).Decode(&doc)
	return doc, err
}

func GenericFindMany[T any](
	coll *mongo.Collection,
	filter bson.D,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	if filter == nil {
		filter = bson.D{{}}
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetLimit(nPerPage)

	paginFilter := bson.D{{"_id", bson.D{{"$lt", startValue}}}}
	andFilter := bson.D{{"$and", bson.A{paginFilter, filter}}}
	if startValue.IsZero() {
		andFilter = filter
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.Find(ctx, andFilter, findOptions)
}

type Service struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Price       uint               `json:"price" bson:"price"`
	Duration    uint               `json:"duration" bson:"duration"`
	Description string             `json:"description" bson:"description"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
}

func (service *Service) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("services")
	return genericInsertOne(coll, service)
}

type Description struct {
	CompanyID   primitive.ObjectID `json:"-" bson:"_id"`
	Description string             `json:"description" bson:"description"`
}

func (description *Description) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("descriptions")
	return genericInsertOne(coll, description)
}

type Company struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Type             string             `json:"type" bson:"type"`
	Localisation     string             `json:"localisation" bson:"localisation"`
	ShortDescription string             `json:"short_description" bson:"short_description"`
}

func (company *Company) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("companies")
	return genericInsertOne(coll, company)
}

// Combined representation of Company and service documents.
// It is used for easier insertion to database
type CompanyCombRepr struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Localisation     string `json:"localisation"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
}

func (companyCombRepr *CompanyCombRepr) InsertCombRepr(db *mongo.Database) (*mongo.InsertOneResult, error) {
	company := Company{
		Name:             companyCombRepr.Name,
		Type:             companyCombRepr.Type,
		Localisation:     companyCombRepr.Localisation,
		ShortDescription: companyCombRepr.ShortDescription,
	}
	result, err := company.InsertOne(db)
	if err != nil {
		return result, err
	}
	description := Description{
		CompanyID:   result.InsertedID.(primitive.ObjectID),
		Description: companyCombRepr.Description,
	}
	return description.InsertOne(db)
}

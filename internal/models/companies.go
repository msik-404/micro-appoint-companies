package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (company *Company) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("companies")
	return GenericInsertOne(coll, company)
}

func (company *Company) UpdateOne(db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection("companies")
	filter := bson.D{{"_id", company.ID}}
	return GenericUpdateOne(coll, filter, company)
}

// Combined representation of Company and service documents.
// It is used for easier insertion to database
type CompanyCombRepr struct {
	ID               primitive.ObjectID `json:"_id"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Localisation     string             `json:"localisation"`
	ShortDescription string             `json:"short_description"`
	Description      string             `json:"description"`
}

func (companyCombRepr *CompanyCombRepr) InsertCombRepr(db *mongo.Database) ([]*mongo.InsertOneResult, error) {
	company := Company{
		Name:             companyCombRepr.Name,
		Type:             companyCombRepr.Type,
		Localisation:     companyCombRepr.Localisation,
		ShortDescription: companyCombRepr.ShortDescription,
	}
	var results []*mongo.InsertOneResult
	result, err := company.InsertOne(db)
	if err != nil {
		return nil, err
	}
	results = append(results, result)
	description := Description{
		CompanyID:   result.InsertedID.(primitive.ObjectID),
		Description: companyCombRepr.Description,
	}
	result, err = description.InsertOne(db)
	if err != nil {
		return nil, err
	}
	results = append(results, result)
	return results, nil
}

func (companyCombRepr *CompanyCombRepr) UpdateCombRepr(db *mongo.Database) ([]*mongo.UpdateResult, error) {
	company := Company{
		ID:               companyCombRepr.ID,
		Name:             companyCombRepr.Name,
		Type:             companyCombRepr.Type,
		Localisation:     companyCombRepr.Localisation,
		ShortDescription: companyCombRepr.ShortDescription,
	}
	var results []*mongo.UpdateResult
	result, err := company.UpdateOne(db)
	results = append(results, result)
	if err != nil {
		return nil, err
	}
	description := Description{
		CompanyID:   companyCombRepr.ID,
		Description: companyCombRepr.Description,
	}
	result, err = description.UpdateOne(db)
	if err != nil {
		return nil, err
	}
	results = append(results, result)
	return results, nil
}

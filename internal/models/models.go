package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func toBsonRemoveEmpty(item any) (doc *bson.M, err error) {
	data, err := bson.Marshal(item)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

type MongoModel interface {
	InsertOne(*mongo.Database) (*mongo.InsertOneResult, error)
	UpdateOne(*mongo.Database) (*mongo.UpdateResult, error)
}

func genericInsertOne[T MongoModel](coll *mongo.Collection, item T) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.InsertOne(ctx, item)
}

func GenericFindOne[T any](coll *mongo.Collection, filter bson.D) (doc T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = coll.FindOne(ctx, filter).Decode(&doc)
	return
}

func GenericUpdateOne[T MongoModel](coll *mongo.Collection, filter bson.D, item T) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	replacement, err := toBsonRemoveEmpty(item)
	delete(*replacement, "_id")
	if err != nil {
		return nil, err
	}
    // empty replacement is skipped
	if len(*replacement) == 0 {
        return nil, nil
	}
    return coll.UpdateOne(ctx, filter, bson.M{"$set": *replacement})
}

func GenericDeleteOne(coll *mongo.Collection, filter bson.D) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.DeleteOne(ctx, filter)
}

func GenericDeleteMany(coll *mongo.Collection, filter bson.D) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.DeleteMany(ctx, filter)
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

	paginFilter := bson.M{"_id": bson.M{"$lt": startValue}}
	andFilter := bson.D{{"$and", bson.A{paginFilter, filter}}}
	if startValue.IsZero() {
		andFilter = filter
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.Find(ctx, andFilter, findOptions)
}

type Service struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Price       uint               `json:"price" bson:"price,omitempty"`
	Duration    uint               `json:"duration" bson:"duration,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id,omitempty"`
}

func (service *Service) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("services")
	return genericInsertOne(coll, service)
}

func (service *Service) UpdateOne(db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection("services")
	filter := bson.D{{"_id", service.ID}}
	return GenericUpdateOne(coll, filter, service)
}

type Description struct {
	CompanyID   primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
}

func (description *Description) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("descriptions")
	return genericInsertOne(coll, description)
}

func (description *Description) UpdateOne(db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection("descriptions")
	filter := bson.D{{"_id", description.CompanyID}}
	return GenericUpdateOne(coll, filter, description)
}

type Company struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name,omitempty"`
	Type             string             `json:"type" bson:"type,omitempty"`
	Localisation     string             `json:"localisation" bson:"localisation,omitempty"`
	ShortDescription string             `json:"short_description" bson:"short_description,omitempty"`
}

func (company *Company) InsertOne(db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection("companies")
	return genericInsertOne(coll, company)
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

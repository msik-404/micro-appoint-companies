package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Price       int                `json:"price" bson:"price,omitempty"`
	Duration    int                `json:"duration" bson:"duration,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
}

type Company struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name,omitempty"`
	Type             string             `json:"type" bson:"type,omitempty"`
	Localisation     string             `json:"localisation" bson:"localisation,omitempty"`
	ShortDescription string             `json:"short_description" bson:"short_description,omitempty"`
	LongDescription  string             `json:"long_description" bson:"long_description,omitempty"`
	Services         []Service          `json:"services" bson:"services,omitempty"`
}

func (company *Company) InsertOne(
	db *mongo.Database,
) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("companies")
	return coll.InsertOne(ctx, company)
}

type CompanyUpdate struct {
	Name             string `json:"name" bson:"name,omitempty"`
	Type             string `json:"type" bson:"type,omitempty"`
	Localisation     string `json:"localisation" bson:"localisation,omitempty"`
	ShortDescription string `json:"short_description" bson:"short_description,omitempty"`
	LongDescription  string `json:"long_description" bson:"long_description,omitempty"`
}

func (companyUpdate *CompanyUpdate) UpdateOne(
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("companies")
	filter := bson.M{"_id": companyID}
	update := bson.M{"$set": companyUpdate}
	return coll.UpdateOne(ctx, filter, update)
}

func DeleteOneCompany(
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("companies")
	filter := bson.M{"_id": companyID}
	return coll.DeleteOne(ctx, filter)
}

func FindOneCompany(
	db *mongo.Database,
	companyID primitive.ObjectID,
) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOne()
	opts.SetProjection(bson.D{
		{"_id", 0},
		{"short_description", 0},
		{"services", bson.M{"$slice": 10}},
	})

	coll := db.Collection("companies")
	filter := bson.M{"_id": companyID}
	return coll.FindOne(ctx, filter, opts)
}

func FindManyCompanies(
	db *mongo.Database,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.M{"_id": -1})
	opts.SetLimit(nPerPage)
	opts.SetProjection(bson.D{
		{"long_description", 0},
		{"services", 0},
	})

	filter := bson.M{}
	if !startValue.IsZero() {
		filter = bson.M{"_id": bson.M{"$lt": startValue}}
	}
	coll := db.Collection("companies")
	return coll.Find(ctx, filter, opts)
}

func (service *Service) InsertOne(
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service.ID = primitive.NewObjectID()

	coll := db.Collection("companies")
	filter := bson.M{"_id": companyID}
	update := bson.M{"$push": bson.M{"services": service}}
	return coll.UpdateOne(ctx, filter, update)
}

func (serviceUpdate *Service) UpdateOne(
	db *mongo.Database,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("companies")
	filter := bson.M{"services._id": serviceID}
    serviceUpdate.ID = serviceID
    update := bson.M{"$set": bson.M{"services.$": serviceUpdate}}
	return coll.UpdateOne(ctx, filter, update)
}

func DeleteOneService(
	db *mongo.Database,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.Collection("companies")
	filter := bson.M{}
	update := bson.M{"$pull": bson.M{"services": bson.M{"_id": serviceID}}}
	return coll.UpdateOne(ctx, filter, update)
}

func FindManyServices(
	db *mongo.Database,
	companyID primitive.ObjectID,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.M{"_id": companyID}}}
	projectionStage := bson.D{{"$project", bson.M{"_id": 0, "services": 1}}}
	unwindStage := bson.D{{"$unwind", "$services"}}
	rootStage := bson.D{{"$replaceRoot", bson.M{"newRoot": "$services"}}}
	sortStage := bson.D{{"$sort", bson.M{"_id": -1}}}
	limitStage := bson.D{{"$limit", nPerPage}}

	pipeline := mongo.Pipeline{
		matchStage,
		projectionStage,
		unwindStage,
		rootStage,
		sortStage,
		limitStage,
	}
	if !startValue.IsZero() {
		startValueStage := bson.D{
			{"$match", bson.M{"_id": bson.M{"$lt": startValue}}},
		}
		pipeline = mongo.Pipeline{
			matchStage,
			projectionStage,
			unwindStage,
			rootStage,
			sortStage,
			startValueStage,
			limitStage,
		}
	}
	coll := db.Collection("companies")
	return coll.Aggregate(ctx, pipeline)
}

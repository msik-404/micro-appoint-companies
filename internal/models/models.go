package models

import (
	"context"
	"time"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" binding:"max=30"`
	Price       int32              `bson:"price,omitempty" binding:"max=1000000"`
	Duration    int32              `bson:"duration,omitempty" binding:"max=480"`
	Description string             `bson:"description,omitempty" binding:"max=300"`
}

type Company struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name,omitempty" binding:"max=30"`
	Type             string             `bson:"type,omitempty" binding:"max=30"`
	Localisation     string             `bson:"localisation,omitempty" binding:"max=60"`
	ShortDescription string             `bson:"short_description,omitempty" binding:"max=150"`
	LongDescription  string             `bson:"long_description,omitempty" binding:"max=300"`
	Services         []Service          `bson:"services,omitempty"`
}

func (company *Company) InsertOne(
    ctx context.Context,
	db *mongo.Database,
) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Collection(database.CollName)
	return coll.InsertOne(ctx, company)
}

type CompanyUpdate struct {
	Name             string `bson:"name,omitempty" binding:"max=30"`
	Type             string `bson:"type,omitempty" binding:"max=30"`
	Localisation     string `bson:"localisation,omitempty" binding:"max=60"`
	ShortDescription string `bson:"short_description,omitempty" binding:"max=150"`
	LongDescription  string `bson:"long_description,omitempty" binding:"max=300"`
}

func (companyUpdate *CompanyUpdate) UpdateOne(
    ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Collection(database.CollName)
	update := bson.M{"$set": companyUpdate}
	return coll.UpdateByID(ctx, companyID, update)
}

func DeleteOneCompany(
    ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Collection(database.CollName)
	filter := bson.M{"_id": companyID}
	return coll.DeleteOne(ctx, filter)
}

func FindOneCompany(
    ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.FindOne()
	opts.SetProjection(bson.D{
		{Key: "_id", Value: 0},
		{Key: "short_description", Value: 0},
		{Key: "services", Value: bson.M{"$slice": 10}},
	})

	coll := db.Collection(database.CollName)
	filter := bson.M{"_id": companyID}
	return coll.FindOne(ctx, filter, opts)
}

func FindManyCompanies(
    ctx context.Context,
	db *mongo.Database,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.M{"_id": -1})
	opts.SetLimit(nPerPage)
	opts.SetProjection(bson.D{
		{Key: "long_description", Value: 0},
		{Key: "services", Value: 0},
	})

	filter := bson.M{}
	if !startValue.IsZero() {
		filter = bson.M{"_id": bson.M{"$lt": startValue}}
	}
	coll := db.Collection(database.CollName)
	return coll.Find(ctx, filter, opts)
}

func (service *Service) InsertOne(
    ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	service.ID = primitive.NewObjectID()

	coll := db.Collection(database.CollName)
	update := bson.M{"$push": bson.M{"services": service}}
	return coll.UpdateByID(ctx, companyID, update)
}

func (serviceUpdate *Service) UpdateOne(
    ctx context.Context,
	db *mongo.Database,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Collection(database.CollName)
	filter := bson.M{"services._id": serviceID}
	serviceUpdate.ID = serviceID
	update := bson.M{"$set": bson.M{"services.$": serviceUpdate}}
	return coll.UpdateOne(ctx, filter, update)
}

func DeleteOneService(
    ctx context.Context,
	db *mongo.Database,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := db.Collection(database.CollName)
    filter := bson.M{"services._id": serviceID}
	update := bson.M{"$pull": bson.M{"services": bson.M{"_id": serviceID}}}
	return coll.UpdateOne(ctx, filter, update)
}

func FindManyServices(
    ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.M{"_id": companyID}}}
	projectionStage := bson.D{{Key: "$project", Value: bson.M{"_id": 0, "services": 1}}}
	unwindStage := bson.D{{Key: "$unwind", Value: "$services"}}
	rootStage := bson.D{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$services"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: nPerPage}}

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
			{Key: "$match", Value: bson.M{"_id": bson.M{"$lt": startValue}}},
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
	coll := db.Collection(database.CollName)
	return coll.Aggregate(ctx, pipeline)
}

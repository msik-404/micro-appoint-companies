package models

import (
	"context"
	"fmt"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	ID          primitive.ObjectID `bson:"service_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Price       int32              `bson:"price,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
	Description string             `bson:"description,omitempty"`
}

type Company struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name,omitempty"`
	Type             string             `bson:"type,omitempty"`
	Localisation     string             `bson:"localisation,omitempty"`
	ShortDescription string             `bson:"short_description,omitempty"`
	LongDescription  string             `bson:"long_description,omitempty"`
	Services         []Service          `bson:"services,omitempty"`
}

func (company *Company) InsertOne(
	ctx context.Context,
	db *mongo.Database,
) (*mongo.InsertOneResult, error) {
	coll := db.Collection(database.CollName)
	return coll.InsertOne(ctx, company)
}

type CompanyUpdate struct {
	Name             *string   `bson:"name,omitempty"`
	Type             *string   `bson:"type,omitempty"`
	Localisation     *string   `bson:"localisation,omitempty"`
	ShortDescription *string   `bson:"short_description,omitempty"`
	LongDescription  *string   `bson:"long_description,omitempty"`
	Services         []Service `bson:"services,omitempty"`
}

func (companyUpdate *CompanyUpdate) UpdateOne(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	coll := db.Collection(database.CollName)
	update := bson.M{"$set": companyUpdate}
	return coll.UpdateByID(ctx, companyID, update)
}

func DeleteOneCompany(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.DeleteResult, error) {
	coll := db.Collection(database.CollName)
	filter := bson.M{"_id": companyID}
	return coll.DeleteOne(ctx, filter)
}

func FindOneCompany(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) *mongo.SingleResult {
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

func FindManyCompaniesByIds(
	ctx context.Context,
	db *mongo.Database,
    companyIDS []primitive.ObjectID,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	opts := options.Find()
	opts.SetSort(bson.M{"_id": -1})
	opts.SetLimit(nPerPage)
	opts.SetProjection(bson.D{
		{Key: "long_description", Value: 0},
		{Key: "services", Value: 0},
	})

    filter := bson.M{"_id": bson.M{"$in": bson.A{companyIDS}}}
	if !startValue.IsZero() {
        filter = bson.M{"$and": bson.A{
            filter,
            bson.M{"_id": bson.M{"$lt": startValue}},
        }}
	}
	coll := db.Collection(database.CollName)
	return coll.Find(ctx, filter, opts)
}

func (service *Service) InsertOne(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	service.ID = primitive.NewObjectID()

	coll := db.Collection(database.CollName)
	update := bson.M{"$push": bson.M{"services": service}}
	return coll.UpdateByID(ctx, companyID, update)
}

func toBsonRemoveEmpty(value any) (doc *bson.M, err error) {
	data, err := bson.Marshal(value)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func getUpdateTerms(updateMap *bson.M) bson.M {
	updateTerms := bson.M{}
	for key, value := range *updateMap {
		if key != "_id" {
			key := fmt.Sprintf("services.$.%s", key)
			updateTerms[key] = value
		}
	}
	return updateTerms
}

type ServiceUpdate struct {
	Name        *string `bson:"name,omitempty"`
	Price       *int32  `bson:"price,omitempty"`
	Duration    *int32  `bson:"duration,omitempty"`
	Description *string `bson:"description,omitempty"`
}

func (serviceUpdate *ServiceUpdate) UpdateOne(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	// this function will erease nil fields,
	// so that unwanted fields will not be set to empty
	updateMap, err := toBsonRemoveEmpty(*serviceUpdate)
	if err != nil {
		return nil, err
	}
	updateTerms := getUpdateTerms(updateMap)

	coll := db.Collection(database.CollName)
	filter := bson.D{
		{Key: "_id", Value: companyID},
		{Key: "services.service_id", Value: serviceID},
	}
	update := bson.M{"$set": updateTerms}
	return coll.UpdateOne(ctx, filter, update)
}

func DeleteOneService(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
	serviceID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	coll := db.Collection(database.CollName)
	filter := bson.M{"$and": bson.A{
		bson.M{"_id": companyID},
		bson.M{"services.service_id": serviceID},
	}}
	update := bson.M{
		"$pull": bson.M{"services": bson.M{"service_id": serviceID}},
	}
	return coll.UpdateOne(ctx, filter, update)
}

func FindManyServices(
	ctx context.Context,
	db *mongo.Database,
	companyID primitive.ObjectID,
	startValue primitive.ObjectID,
	nPerPage int64,
) (*mongo.Cursor, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"_id": companyID}}}
	projectionStage := bson.D{{Key: "$project", Value: bson.M{"_id": 0, "services": 1}}}
	unwindStage := bson.D{{Key: "$unwind", Value: "$services"}}
	rootStage := bson.D{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$services"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"service_id": -1}}}
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
			{Key: "$match", Value: bson.M{"service_id": bson.M{"$lt": startValue}}},
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

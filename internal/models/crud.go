package models

import (
	"context"
	"fmt"
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

func GenericInsertOne[T MongoModel](coll *mongo.Collection, item T) (*mongo.InsertOneResult, error) {
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
    fmt.Printf("%+v", filter)
    fmt.Printf("%+v", *replacement)
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

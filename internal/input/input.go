package input

import (
	"context"
	"strconv"
	"time"
    "errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findPagin returns paginated results.
// It first finds all documents with id lower than startValue, then sorts and limits.
func FindPagin(coll *mongo.Collection, inputFilter bson.D, startValue primitive.ObjectID, nPerPage uint) (*mongo.Cursor, error) {
	if inputFilter == nil {
		inputFilter = bson.D{{}}
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetLimit(int64(nPerPage))

	paginFilter := bson.D{{"_id", bson.D{{"$lt", startValue}}}}
	filter := bson.D{{"$and", bson.A{paginFilter, inputFilter}}}
	if startValue.IsZero() {
		filter = inputFilter
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return coll.Find(ctx, filter, findOptions)
}

func GetObjectId(hexObjectIdSting string) (primitive.ObjectID, error) {
	id := primitive.NilObjectID
	if hexObjectIdSting != "" {
		return primitive.ObjectIDFromHex(hexObjectIdSting)
	}
	return id, nil
}

func GetStartValue(c *gin.Context) (primitive.ObjectID, error) {
	query := c.DefaultQuery("startValue", "")
	return GetObjectId(query)
}

func GetNPerPageValue(c *gin.Context) (uint, error) {
	query := c.DefaultQuery("nPerPage", "100")
	nPerPage, err := strconv.Atoi(query)
	if err != nil {
		return uint(nPerPage), err
	}
	if nPerPage < 0 {
		return uint(nPerPage), errors.New("nPerPage should be positive number")
	}
	return uint(nPerPage), nil
}

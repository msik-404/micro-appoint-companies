package pagination

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msik-404/micro-appoint-companies/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findPagin returns paginated results.
// It first finds all documents with id lower than startValue, then sorts and limits.
func FindPagin(coll *mongo.Collection, inputFilter bson.D, startValue primitive.ObjectID, nPerPage uint) (*mongo.Cursor, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetLimit(int64(nPerPage))

	paginFilter := bson.D{{"_id", bson.D{{"$lt", startValue}}}}
	filter := bson.D{{"$and", bson.A{paginFilter, inputFilter}}}
	if startValue.IsZero() {
		filter = inputFilter
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return coll.Find(ctx, filter, findOptions)
}

func GetStartValue(c *gin.Context) (primitive.ObjectID, error) {
	query := c.DefaultQuery("startValue", "")
	return utils.GetObjectId(query)
}

type NPerPageParseErorr struct{}

func (e *NPerPageParseErorr) Error() string {
	return "nPerPage should be positive number"
}

func GetNPerPageValue(c *gin.Context) (uint, error) {
	query := c.DefaultQuery("nPerPage", "100")
	nPerPage, err := strconv.Atoi(query)
	if err != nil {
		return uint(nPerPage), err
	}
	if nPerPage < 0 {
		return uint(nPerPage), &NPerPageParseErorr{}
	}
	return uint(nPerPage), nil
}

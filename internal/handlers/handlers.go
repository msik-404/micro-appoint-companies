package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/msik-404/micro-appoint-companies/internal/models"
)

// findPagin returns paginated results.
// It first finds all documents with id lower than startValue, then sorts and limits.
func findPagin(coll *mongo.Collection, startValue primitive.ObjectID, nPerPage uint) (*mongo.Cursor, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetLimit(int64(nPerPage))
	filter := bson.D{{"_id", bson.D{{"$lt", startValue}}}}
	if startValue.IsZero() {
		filter = bson.D{{}}
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return coll.Find(ctx, filter, findOptions)
}

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		query := c.DefaultQuery("startValue", "")
		startValue := primitive.NilObjectID
		if query != "" {
			providedStartValue, err := primitive.ObjectIDFromHex(query)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			startValue = providedStartValue
		}
		query = c.DefaultQuery("nPerPage", "100")
		nPerPage, err := strconv.Atoi(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if nPerPage < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nPerPage should be positive number"})
			return
		}
		coll := db.Collection("companies")
		cursor, err := findPagin(coll, startValue, uint(nPerPage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var companies []models.Company
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := cursor.All(ctx, &companies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, companies)
	}
	return gin.HandlerFunc(fn)
}

func GetServicesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

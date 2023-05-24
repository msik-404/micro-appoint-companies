package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/models"
)

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

func GetNPerPageValue(c *gin.Context) (int64, error) {
	query := c.DefaultQuery("nPerPage", "100")
	nPerPage, err := strconv.Atoi(query)
	if err != nil {
		return int64(nPerPage), err
	}
	if nPerPage < 0 {
		return int64(nPerPage), errors.New("nPerPage should be positive number")
	}
	return int64(nPerPage), nil
}

func PaginHandler[T any](c *gin.Context, coll *mongo.Collection, filter bson.D) ([]T, error) {
	// parse variables for pagination
	startValue, err := GetStartValue(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	nPerPage, err := GetNPerPageValue(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	// get records based on pagination variables
	cursor, err := models.GenericFindMany[T](coll, filter, startValue, nPerPage)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	// transfom cursor into slice of results
	var results []T
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cursor.All(ctx, &results); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	if len(results) == 0 {
		err = errors.New("No documents in the result")
		c.AbortWithError(http.StatusNotFound, err)
		return nil, err
	}
	return results, err
}

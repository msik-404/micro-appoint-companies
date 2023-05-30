package middleware

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

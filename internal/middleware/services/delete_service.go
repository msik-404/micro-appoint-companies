package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func DeleteServiceEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serviceID, err := middleware.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("services")
		filter := bson.D{{"_id", serviceID}}
		result, err := models.GenericDeleteOne(coll, filter)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

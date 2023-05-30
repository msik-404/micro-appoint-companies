package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func UpdateServiceEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serviceID, err := middleware.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var serviceUpdate models.Service
		if err := c.BindJSON(&serviceUpdate); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		result, err := serviceUpdate.UpdateOne(db, serviceID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

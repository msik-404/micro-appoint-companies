package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		var newService models.Service
		if err := c.BindJSON(&newService); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		newService.ID = serviceID
		// Update of service is not allowed to change company to which it belongs.
		// Maybe admin should be allowed to do that.
		newService.CompanyID = primitive.NilObjectID
		result, err := newService.UpdateOne(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

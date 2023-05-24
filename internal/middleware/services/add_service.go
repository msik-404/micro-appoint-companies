package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func AddServiceEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newService models.Service
		if err := c.BindJSON(&newService); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		result, err := newService.InsertOne(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

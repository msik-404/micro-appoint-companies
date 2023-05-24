package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func GetServicesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := middleware.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("services")
		filter := bson.D{{"company_id", companyID}}
		services, err := middleware.PaginHandler[models.Service](c, coll, filter)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, services)
	}
	return gin.HandlerFunc(fn)
}

package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		startValue, err := middleware.GetStartValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		nPerPage, err := middleware.GetNPerPageValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		cursor, err := models.FindManyServices(db, companyID, startValue, nPerPage)
		var services []models.Service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cursor.All(ctx, &services); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if len(services) == 0 {
			err = errors.New("No documents in the result")
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, services)
	}
	return gin.HandlerFunc(fn)
}

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/models"
	"github.com/msik-404/micro-appoint-companies/internal/pagination"
)

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		startValue, err := pagination.GetStartValue(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		coll := db.Collection("companies")
        nPerPage, err := pagination.GetNPerPageValue(c)
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
		cursor, err := pagination.FindPagin(coll, startValue, uint(nPerPage))
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

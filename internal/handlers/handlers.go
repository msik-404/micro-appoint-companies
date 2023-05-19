package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/models"
	"github.com/msik-404/micro-appoint-companies/internal/pagination"
	"github.com/msik-404/micro-appoint-companies/internal/utils"
)

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
        // parse variables for pagination
		startValue, err := pagination.GetStartValue(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
        nPerPage, err := pagination.GetNPerPageValue(c)
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        // get records based on pagination variables
        coll := db.Collection("companies")
        filter := bson.D{{}}
		cursor, err := pagination.FindPagin(coll, filter, startValue, uint(nPerPage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
        // transfom cursor into slice of companies
		var companies []bson.M
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := cursor.All(ctx, &companies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, companies)
	}
	return gin.HandlerFunc(fn)
}

func GetCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
    fn := func(c *gin.Context) {
        companyID, err := utils.GetObjectId(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
        coll := db.Collection("descriptions")
        filter := bson.D{{"_id", companyID}}
	    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
        var descDoc models.Description
        err = coll.FindOne(ctx, filter).Decode(&descDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
        c.JSON(http.StatusOK, descDoc)
    }
    return gin.HandlerFunc(fn)
}

func GetServicesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
        companyID, err := utils.GetObjectId(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
        // parse variables for pagination
		startValue, err := pagination.GetStartValue(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
        nPerPage, err := pagination.GetNPerPageValue(c)
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        coll := db.Collection("services")
        filter := bson.D{{"company_id", companyID}}
        cursor, err := pagination.FindPagin(coll, filter, startValue, uint(nPerPage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
        // transfom cursor into slice of services
		var services []models.Service
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := cursor.All(ctx, &services); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)
	}
	return gin.HandlerFunc(fn)
}

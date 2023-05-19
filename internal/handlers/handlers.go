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
		coll := db.Collection("companies")
        nPerPage, err := pagination.GetNPerPageValue(c)
        if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        // get records based on pagination variables
		cursor, err := pagination.FindPagin(coll, startValue, uint(nPerPage))
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

	}
	return gin.HandlerFunc(fn)
}

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/input"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// parse variables for pagination
		startValue, err := input.GetStartValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		nPerPage, err := input.GetNPerPageValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// get records based on pagination variables
		coll := db.Collection("companies")
		cursor, err := input.FindPagin(coll, nil, startValue, uint(nPerPage))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// transfom cursor into slice of companies
		var companies []bson.M
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cursor.All(ctx, &companies); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, companies)
	}
	return gin.HandlerFunc(fn)
}

func GetCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := input.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("descriptions")
		filter := bson.D{{"_id", companyID}}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var descDoc models.Description
		err = coll.FindOne(ctx, filter).Decode(&descDoc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, descDoc)
	}
	return gin.HandlerFunc(fn)
}

func GetServicesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := input.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// parse variables for pagination
		startValue, err := input.GetStartValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		nPerPage, err := input.GetNPerPageValue(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("services")
		filter := bson.D{{"company_id", companyID}}
		cursor, err := input.FindPagin(coll, filter, startValue, uint(nPerPage))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// transfom cursor into slice of services
		var services []models.Service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cursor.All(ctx, &services); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, services)
	}
	return gin.HandlerFunc(fn)
}

func PostCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newCompany models.CompanyPost
		if err := c.BindJSON(&newCompany); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		result, err := newCompany.Insert(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

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

func paginHandler[T any](c *gin.Context, coll *mongo.Collection, filter bson.D) ([]T, error) {
	// parse variables for pagination
	startValue, err := input.GetStartValue(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	nPerPage, err := input.GetNPerPageValue(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	// get records based on pagination variables
	cursor, err := models.GenericFindMany[T](coll, filter, startValue, nPerPage)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	// transfom cursor into slice of results
	var results []T
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cursor.All(ctx, &results); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	return results, nil
}

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		coll := db.Collection("companies")
		companies, err := paginHandler[models.Company](c, coll, nil)
		if err != nil {
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
		descDoc, err := models.GenericFindOne[models.Description](coll, filter)
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
		coll := db.Collection("services")
		filter := bson.D{{"company_id", companyID}}
		services, err := paginHandler[models.Service](c, coll, filter)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, services)
	}
	return gin.HandlerFunc(fn)
}

func PostCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newCompany models.CompanyCombRepr
		if err := c.BindJSON(&newCompany); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		result, err := newCompany.InsertCombRepr(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

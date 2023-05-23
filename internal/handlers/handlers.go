package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/input"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

// if rseource not found 404 should be returned

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
	if len(results) == 0 {
		err = errors.New("No documents in the result")
		c.AbortWithError(http.StatusNotFound, err)
		return nil, err
	}
	return results, err
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
			if err == mongo.ErrNoDocuments {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
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

func AddCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newCompany models.CompanyCombRepr
		if err := c.BindJSON(&newCompany); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		results, err := newCompany.InsertCombRepr(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, results)
	}
	return gin.HandlerFunc(fn)
}

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

func UpdateCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := input.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var newCompany models.CompanyCombRepr
		if err := c.BindJSON(&newCompany); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		newCompany.ID = companyID
		results, err := newCompany.UpdateCombRepr(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, results)
	}
	return gin.HandlerFunc(fn)
}

func UpdateServiceEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serviceID, err := input.GetObjectId(c.Param("id"))
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

// cascade delete
func DeleteCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := input.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("companies")
		filter := bson.D{{"_id", companyID}}
		var results []*mongo.DeleteResult
		result, err := models.GenericDeleteOne(coll, filter)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		results = append(results, result)
		coll = db.Collection("descriptions")
		result, err = models.GenericDeleteOne(coll, filter)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		results = append(results, result)
		coll = db.Collection("services")
		filter = bson.D{{"company_id", companyID}}
		result, err = models.GenericDeleteMany(coll, filter)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		results = append(results, result)
		c.JSON(http.StatusOK, results)
	}
	return gin.HandlerFunc(fn)
}

func DeleteServiceEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		serviceID, err := input.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		coll := db.Collection("services")
		filter := bson.D{{"_id", serviceID}}
		result, err := models.GenericDeleteOne(coll, filter)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

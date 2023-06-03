package companies

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// parse variables for pagination
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
		// get records based on pagination variables
		cursor, err := models.FindManyCompanies(db, startValue, nPerPage)
		// define json structure for output
		type Company struct {
			ID               primitive.ObjectID `json:"_id" bson:"_id"`
			Name             string             `json:"name" bson:"name"`
			Type             string             `json:"type" bson:"type"`
			Localisation     string             `json:"localisation" bson:"localisation"`
			ShortDescription string             `json:"short_description" bson:"short_description"`
		}
		var companies []Company
		// transfom cursor into slice of results
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cursor.All(ctx, &companies); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if len(companies) == 0 {
			err = errors.New("No documents in the result")
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, companies)
	}
	return gin.HandlerFunc(fn)
}

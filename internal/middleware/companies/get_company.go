package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func GetCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := middleware.GetObjectId(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		type Company struct {
            Name            string           `json:"name" bson:"name"`
            Type            string           `json:"type" bson:"type"`
            Localisation    string           `json:"localisation" bson:"localisation"`
            LongDescription string           `json:"long_description" bson:"long_description"`
            Services        []models.Service `json:"services" bson:"services"`
		}
        var company Company

		err = models.FindOneCompany(db, companyID).Decode(&company)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		c.JSON(http.StatusOK, company)
	}
	return gin.HandlerFunc(fn)
}

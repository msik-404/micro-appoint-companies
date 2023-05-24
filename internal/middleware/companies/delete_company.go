package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

// cascade delete
func DeleteCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := middleware.GetObjectId(c.Param("id"))
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

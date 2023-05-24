package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func AddCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newCompany models.Company
		if err := c.BindJSON(&newCompany); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		results, err := newCompany.InsertOne(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, results)
	}
	return gin.HandlerFunc(fn)
}

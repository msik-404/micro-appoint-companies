package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func UpdateCompanyEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		companyID, err := middleware.GetObjectId(c.Param("id"))
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

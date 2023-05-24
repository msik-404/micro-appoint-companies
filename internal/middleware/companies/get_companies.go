package companies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"


	"github.com/msik-404/micro-appoint-companies/internal/middleware"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func GetCompaniesEndPoint(db *mongo.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		coll := db.Collection("companies")
		companies, err := middleware.PaginHandler[models.Company](c, coll, nil)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, companies)
	}
	return gin.HandlerFunc(fn)
}


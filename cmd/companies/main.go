package main

import (
	"github.com/gin-gonic/gin"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"github.com/msik-404/micro-appoint-companies/internal/middleware/companies"
	"github.com/msik-404/micro-appoint-companies/internal/middleware/services"
)

func main() {
	mongoClient, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	db := mongoClient.Database("micro-appoint-companies")
	// _, err = database.CreateDBIndexes(db)
	// if err != nil {
	// 	panic(err)
	// }
	r := gin.Default()

	r.GET("/companies", companies.GetCompaniesEndPoint(db))

	r.GET("/companies/:id", companies.GetCompanyEndPoint(db))
	r.GET("/companies/services/:id", services.GetServicesEndPoint(db))

	r.POST("/companies", companies.AddCompanyEndPoint(db))
    r.POST("/companies/services/:id", services.AddServiceEndPoint(db))

	r.PUT("/companies/:id", companies.UpdateCompanyEndPoint(db))
	r.PUT("/services/:id", services.UpdateServiceEndPoint(db))

	r.DELETE("/companies/:id", companies.DeleteCompanyEndPoint(db))
	r.DELETE("/services/:id", services.DeleteServiceEndPoint(db))

	r.Run() // listen and serve on 0.0.0.0:8080
}

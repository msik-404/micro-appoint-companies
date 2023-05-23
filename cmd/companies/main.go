package main

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"github.com/msik-404/micro-appoint-companies/internal/handlers"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func main() {
	mongoClient, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	db := mongoClient.Database("micro-appoint-companies")
	_, err = database.CreateDBIndexes(db)
	if err != nil {
		panic(err)
	}
	// testInsert(db)

	r := gin.Default()

	r.GET("/companies", handlers.GetCompaniesEndPoint(db))

	r.GET("/companies/:id", handlers.GetCompanyEndPoint(db))
	r.GET("/companies/services/:id", handlers.GetServicesEndPoint(db))

	r.POST("/companies", handlers.AddCompanyEndPoint(db))
	r.POST("/services", handlers.AddServiceEndPoint(db))

	r.PUT("/companies/:id", handlers.UpdateCompanyEndPoint(db))
	r.PUT("/services/:id", handlers.UpdateServiceEndPoint(db))

	r.DELETE("/companies/:id", handlers.DeleteCompanyEndPoint(db))
	r.DELETE("/services/:id", handlers.DeleteServiceEndPoint(db))

	r.Run() // listen and serve on 0.0.0.0:8080
}

func testInsert(db *mongo.Database) {
	collCompany := db.Collection("companies")
	collDesc := db.Collection("descriptions")
	collService := db.Collection("services")
	for i := 0; i < 10; i++ {
		name := "name: " + strconv.Itoa(i)
		companyType := "type: " + strconv.Itoa(i)
		localistaion := "loc: " + strconv.Itoa(i)
		shortDescription := "short description: " + strconv.Itoa(i)
		docCompany := models.Company{
			Name:             name,
			Type:             companyType,
			Localisation:     localistaion,
			ShortDescription: shortDescription,
		}
		result, err := collCompany.InsertOne(context.TODO(), docCompany)
		if err != nil {
			panic(err)
		}
		companyID := result.InsertedID.(primitive.ObjectID)
		description := "description: " + strconv.Itoa(i)
		docDesc := models.Description{
			CompanyID:   companyID,
			Description: description,
		}
		result, err = collDesc.InsertOne(context.TODO(), docDesc)
		if err != nil {
			panic(err)
		}
		for j := 0; j < 15; j++ {
			name := "name: " + strconv.Itoa(j)
			price := uint(i * j)
			description := "description: " + strconv.Itoa(j)
			docService := models.Service{
				Name:        name,
				Price:       price,
				Duration:    60,
				Description: description,
				CompanyID:   companyID,
			}

			result, err = collService.InsertOne(context.TODO(), docService)
			if err != nil {
				panic(err)
			}
		}
	}
}

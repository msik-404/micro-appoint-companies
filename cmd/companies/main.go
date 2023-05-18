package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"github.com/msik-404/micro-appoint-companies/internal/handlers"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

func testInsert(client *mongo.Client) {
	docCompany := models.Company{Name: "jedna", Description: "fajny opis"}
	collCompany := client.Database("micro-appoint-companies").Collection("companies")
	collService := client.Database("micro-appoint-companies").Collection("services")
	for i := 0; i < 10; i++ {
		result, err := collCompany.InsertOne(context.TODO(), docCompany)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

		for j := 0; j < 15; j++ {
			docService := models.Service{
				Name:        "taki",
				Price:       100,
				Duration:    100,
				Description: "dhjfjdsh",
				CompanyID:   result.InsertedID.(primitive.ObjectID),
			}

			result, err = collService.InsertOne(context.TODO(), docService)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
		}
	}
}

func main() {
	mongoClient, err := database.ConnectDB()
    // testInsert(mongoClient)
	if err != nil {
		panic(err)
	}
	db := mongoClient.Database("micro-appoint-companies")

	r := gin.Default()

	r.GET("/companies", handlers.GetCompaniesEndPoint(db))
	// r.GET("/services", handlers.GetServicesEndPoint(db, 0, 100))

	r.Run() // listen and serve on 0.0.0.0:8080
}

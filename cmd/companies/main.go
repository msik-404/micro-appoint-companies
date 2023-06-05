package main

import (
	// "github.com/gin-gonic/gin"

	"fmt"
	"net"

	"github.com/msik-404/micro-appoint-companies/internal/communication"
	"github.com/msik-404/micro-appoint-companies/internal/database"
	"google.golang.org/grpc"
	// "github.com/msik-404/micro-appoint-companies/internal/controllers/companies"
	// "github.com/msik-404/micro-appoint-companies/internal/controllers/services"
)

func main() {
	mongoClient, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	_, err = database.CreateDBIndexes(mongoClient)
	if err != nil {
		panic(err)
	}
    port := 50051
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) 
    if err != nil {
        panic(err)
    }
    s := grpc.NewServer()
    communication.RegisterApiServer(s, &communication.Server{Client: *mongoClient})
    if err := s.Serve(lis); err != nil {
        panic(err)
    }

	// r := gin.Default()

	// r.GET("/companies", companies.GetCompaniesEndPoint(db))

	// r.GET("/companies/:id", companies.GetCompanyEndPoint(db))
	// r.GET("/companies/services/:id", services.GetServicesEndPoint(db))

	// r.POST("/companies", companies.AddCompanyEndPoint(db))
    // r.POST("/companies/services/:id", services.AddServiceEndPoint(db))

	// r.PUT("/companies/:id", companies.UpdateCompanyEndPoint(db))
	// r.PUT("/services/:id", services.UpdateServiceEndPoint(db))

	// r.DELETE("/companies/:id", companies.DeleteCompanyEndPoint(db))
	// r.DELETE("/services/:id", services.DeleteServiceEndPoint(db))

	// r.Run() // listen and serve on 0.0.0.0:8080
}

package main

import (
	"fmt"
	"net"

    "google.golang.org/grpc"

	"github.com/msik-404/micro-appoint-companies/internal/companiespb"
	"github.com/msik-404/micro-appoint-companies/internal/database"
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
    companiespb.RegisterApiServer(s, &companiespb.Server{Client: mongoClient})
    if err := s.Serve(lis); err != nil {
        panic(err)
    }
}

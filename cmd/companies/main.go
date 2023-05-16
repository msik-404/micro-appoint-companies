package main

import (
	"github.com/gin-gonic/gin"
	"github.com/msik-404/micro-appoint-companies/internal/database"
	"net/http"
)

func main() {
	database.ConnectDB()
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}

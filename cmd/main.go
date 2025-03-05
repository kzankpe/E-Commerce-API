package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	//Initializing server
	server = gin.Default()
}

func main() {
	// Group the Api under api
	router := server.Group("/api")

	// Add healthcheck
	router.GET("/healthcheck", func(c *gin.Context) {
		message := "Welcome to the E-Commerce API"
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	//Add Home page
	router.GET("/", func(c *gin.Context) {
		message := "Welcome to the Service. E-Commerce platform with cart and payment gateway integration"
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	// Run the server on a port
	err := server.Run(":8090")
	if err != nil {
		panic(err)
	}
}

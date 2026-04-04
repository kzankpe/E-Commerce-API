package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal/infrastructure/database"
)

var (
	server *gin.Engine
)

func init() {
	//Initializing server
	server = gin.Default()
	datasource := "user=db dbname=db" // move to config
	database.ConnectDB(datasource)

}

func main() {

	log.Println("Starting the server")
	//router := server.Group("/api/v1")

	// Run the server (listen on all interfaces)
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}

}

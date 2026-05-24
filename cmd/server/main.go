package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kzankpe/e-commerce-api/internal/config"

	"github.com/kzankpe/e-commerce-api/pkg/logger"
)

var (
	server *gin.Engine
)

func init() {
	//Initializing server
	server = gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Environment)

	// Connect to database
	// dbConfig := database.Config{
	// 	Host:     cfg.Database.Host,
	// 	Port:     cfg.Database.Port,
	// 	User:     cfg.Database.User,
	// 	Password: cfg.Database.Password,
	// 	DBName:   cfg.Database.Name,
	// 	SSLMode:  cfg.Database.SSLMode,
	// }

	// db, err := database.NewConnection(dbConfig)
	// if err != nil {
	// 	appLogger.Fatal("Failed to connect to database", err)
	// }
	// defer db.Close()

	appLogger.Info("Database connected successfully")

}

func main() {

	log.Println("Starting the server")
	//router := server.Group("/api/v1")
	//handlers := handler.NewAuthHandler(registerUC)
	//r := router.NewRouter(handlers)

	// Run the server (listen on all interfaces)
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}

}

package main

import (
	"os"

	"github.com/sachinggsingh/PDTS/auth-service/config"
	"github.com/sachinggsingh/PDTS/auth-service/internal/api"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"github.com/sachinggsingh/PDTS/auth-service/pkg/db"
)

func main() {
	// Initialize logger and config
	logger := utils.New()
	env := config.GetEnv()

	// Connect to database
	database := db.NewDb()
	if err := database.ConnectDb(); err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		os.Exit(1)
	}
	defer database.Disconnect()

	// Initialize server with all dependencies
	server := api.NewServer(env, logger, database)

	// Setup and get the router with all routes configured
	router := server.APIServer()

	// Start the server
	logger.Info("Starting server on port " + env.PORT)
	router.Run(":" + env.PORT)
}

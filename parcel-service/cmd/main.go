package main

import (
	"fmt"
	"os"

	"github.com/sachinggsingh/PDTS/parcel-service/config"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/api"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
	"github.com/sachinggsingh/PDTS/parcel-service/pkg/db"
)

func main() {

	fmt.Println("Hello from the Parcel Service")

	database := db.NewDatabase()
	if err := database.ConnectToDatabase(); err != nil {
		utils.Log.Error("Error in Connecting to the database")
		os.Exit(1)

	}
	defer database.DisconnectFromDatabase()

	server := api.NewServer(config.LoadENV(), utils.Log, database)
	server.RunServer()
}

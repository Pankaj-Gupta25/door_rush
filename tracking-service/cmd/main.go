package main

import (
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/api"
)

func main() {

	server := api.NewServer()
	server.StartServer()

}

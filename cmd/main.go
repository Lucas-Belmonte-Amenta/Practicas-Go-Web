package main

import (
	"PRACTICAS-GO-WEB/cmd/server"
	"fmt"
	"log"
	"os"
)

func main() {

	port := ":8080"

	cfg := &server.ConfigServer{
		ServerAddress:   port,
		StaticFilesPath: "./docs/db/products.json",
	}

	log.Printf("Server running on port %s", port)
	app := server.NewServer(cfg)

	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

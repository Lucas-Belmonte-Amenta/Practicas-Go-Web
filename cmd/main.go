package main

import (
	"PRACTICAS-GO-WEB/cmd/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Overload()
	if err != nil {
		panic("Error loading .env file")
	}

	token := os.Getenv("Token")
	port := os.Getenv("Port")

	cfg := &server.ConfigServer{
		ServerAddress:   ":" + port,
		StaticFilesPath: "./docs/db/products.json",
	}

	log.Printf("Server running on port %s", port)
	app := server.NewServer(cfg)

	if err := app.Run(token); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

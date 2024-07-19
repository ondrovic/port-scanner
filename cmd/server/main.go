package main

import (
	"log"
	"os"

	"port-scanner/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	router := api.SetupRouter()
	log.Fatal(router.Run(":" + port))
}
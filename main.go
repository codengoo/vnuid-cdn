package main

import (
	"log"
	"os"
	"vnuid-cdn/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New()

	routers.SetupRoutes(app)

	log.Printf("Server started at http://0.0.0.0:%s", port)
	log.Fatal(app.Listen(":" + port))
}

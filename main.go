package main

import (
	"log"
	"os"
	"vnuid_cdn/routers"

	"github.com/gofiber/fiber/v2"
)

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

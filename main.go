package main

import (
	"log"
	"vnuid_cdn/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routers.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}

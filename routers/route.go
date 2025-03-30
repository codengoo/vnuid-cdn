package routers

import (
	"os"
	"vnuid-cdn/middlewares"
	"vnuid-cdn/services"

	"github.com/gofiber/fiber/v2"
)

const (
	UploadDir = "uploads"
)

func SetupRoutes(app *fiber.App) {
	const uploadDir = services.UploadDir

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755) // mkdir -p
	}

	// Group routes for API
	api := app.Group("/api")

	// Upload route with JWT middleware
	api.Post("/upload", middlewares.VerifyToken, services.UploadHandler)
	api.Delete("/delete/:uuid", middlewares.VerifyToken, services.DeleteHandler)
	api.Delete("/delete", middlewares.VerifyToken, services.DeleteMultipleHandler)

	app.Static("/cdn", uploadDir)
}

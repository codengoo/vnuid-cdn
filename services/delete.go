package services

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func DeleteHandler(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing UUID"})
	}

	files, err := filepath.Glob(filepath.Join(UploadDir, uuid+".*"))
	if err != nil || len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File not found"})
	}

	for _, filePath := range files {
		if err := os.Remove(filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete file"})
		}
	}

	return c.JSON(fiber.Map{"message": "File deleted successfully"})
}

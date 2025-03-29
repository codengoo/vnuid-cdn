package services

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type DeleteRequest struct {
	UUIDs []string `json:"uuids"`
}

func DeleteMultipleHandler(c *fiber.Ctx) error {
	var req DeleteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(req.UUIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No UUIDs provided"})
	}

	var deletedFiles []string
	var notFoundFiles []string
	var failedFiles []string

	for _, uuid := range req.UUIDs {
		files, err := filepath.Glob(filepath.Join(UploadDir, uuid+".*"))
		if err != nil || len(files) == 0 {
			notFoundFiles = append(notFoundFiles, uuid)
			continue
		}

		for _, filePath := range files {
			if err := os.Remove(filePath); err != nil {
				failedFiles = append(failedFiles, uuid)
			} else {
				deletedFiles = append(deletedFiles, uuid)
			}
		}
	}

	return c.JSON(fiber.Map{
		"deleted":   deletedFiles,
		"not_found": notFoundFiles,
		"failed":    failedFiles,
	})
}

package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UploadDir         = "uploads"
	maxFileSize int64 = 5 * 1024 * 1024 // 5MB
	maxUploads        = 5
)

var (
	allowedExts = map[string]bool{
		".png": true,
		".jpg": true,
	}
)

func UploadHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error retrieving files"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No files uploaded"})
	} else if len(files) > maxUploads {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot upload more than 5 files at a time"})
	}

	var uploadedFiles []string

	for _, file := range files {
		// Check file size
		if file.Size > maxFileSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds 5MB limit"})
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only PNG, JPEG file are allowed"})
		}

		// rename the file
		newFileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join(UploadDir, newFileName)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		url := fmt.Sprintf("/cdn/%s", newFileName)
		uploadedFiles = append(uploadedFiles, url)
	}
	return c.JSON(fiber.Map{"urls": uploadedFiles})
}

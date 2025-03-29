package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	uploadDir         = "uploads"
	maxFileSize int64 = 5 * 1024 * 1024 // 5MB
	maxUploads        = 5
	tokenHeader       = "Authorization"
	jwtSecret         = "your-secret-key"
)

var (
	allowedExts = map[string]bool{
		".png": true,
		".jpg": true,
	}
)

func verifyToken(c *fiber.Ctx) error {
	tokenString := c.Get(tokenHeader)

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	return c.Next()
}

func handleUpload(c *fiber.Ctx) error {
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
		filePath := filepath.Join(uploadDir, newFileName)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		url := fmt.Sprintf("/cdn/%s", newFileName)
		uploadedFiles = append(uploadedFiles, url)
	}
	return c.JSON(fiber.Map{"urls": uploadedFiles})
}

func main() {

	app := fiber.New()

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755) // mkdir -p
	}

	app.Post("/upload", verifyToken, handleUpload)

	app.Static("/cdn", uploadDir)
	log.Fatal(app.Listen(":8080"))
}

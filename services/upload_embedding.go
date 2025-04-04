package services

import (
	"encoding/json"
	"fmt"
	"os"
	"vnuid-cdn/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadEmbeddingRequest struct {
	Embedding []float64 `json:"embedding" validate:"required"`
	Uid       string    `json:"uid" validate:"required"`
}

const EmbeddingDir = "embedding"

func UploadEmbedding(ctx *fiber.Ctx) error {
	var data UploadEmbeddingRequest

	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if msgs := utils.Validate(data); msgs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid args", "msgs": msgs})
	}

	userDir := EmbeddingDir + "/" + data.Uid
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		os.MkdirAll(userDir, 0755) // mkdir -p
	}

	filename := uuid.New().String() + ".json"
	filepath := userDir + "/" + filename
	file, err := os.Create(filepath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while creating file"})
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(data); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while writing file"})
	}

	url := fmt.Sprintf("/emb/%s/%s", data.Uid, filename)
	return ctx.JSON(fiber.Map{"filename": url})
}

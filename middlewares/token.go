package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenHeader = "Authorization"
)

var jwtSecret string

func init() {
	jwtSecret = os.Getenv("JWT_KEY")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}
}

func VerifyToken(c *fiber.Ctx) error {
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

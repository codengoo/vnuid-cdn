package middlewares

import "github.com/gofiber/fiber/v2"

func LongtimeCache(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public, max-age=2592000, immutable") // ~1 month
	return c.Next()
}

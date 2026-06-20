package handlers

import "github.com/gofiber/fiber/v2"

func respondError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message,
	})
}

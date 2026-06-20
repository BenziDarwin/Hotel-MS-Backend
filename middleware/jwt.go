package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized",
			})
		},
	})
}

func UserIDFromToken(c *fiber.Ctx) string {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	if sub, ok := claims["sub"].(string); ok {
		return sub
	}

	return ""
}

func RoleFromToken(c *fiber.Ctx) string {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	if role, ok := claims["role"].(string); ok {
		return strings.TrimSpace(role)
	}

	return ""
}

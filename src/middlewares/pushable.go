package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func WithPushable() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Method() == "POST" || ctx.Method() == "PUT" {
			return ctx.Next()
		}

		return ctx.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid http request method",
		})
	}
}

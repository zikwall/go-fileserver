package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func WithFilename() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		paths := strings.Split(ctx.Path(), "/")

		ctx.Locals("filename", paths[len(paths)-1])

		return ctx.Next()
	}
}

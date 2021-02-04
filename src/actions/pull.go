package actions

import (
	"github.com/gofiber/fiber/v2"
)

func (a *ActionProvider) PullFile(ctx *fiber.Ctx) error {
	filepath, _, code, err := a.GetFilepath(ctx)

	if err != nil {
		return ctx.Status(code).JSON(Response{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.SendFile(filepath)
}

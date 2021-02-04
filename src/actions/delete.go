package actions

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func (a *ActionProvider) DeleteFile(ctx *fiber.Ctx) error {
	filepath, filename, code, err := a.GetFilepath(ctx)

	if err != nil {
		return ctx.Status(code).JSON(Response{
			Status:  false,
			Message: err.Error(),
		})
	}

	if err := os.Remove(filepath); err != nil {
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: fmt.Sprintf("Can't delete file: %s", filename),
		})
	}

	return ctx.JSON(Response{
		Status:  true,
		Message: "File successfully deleted!",
	})
}

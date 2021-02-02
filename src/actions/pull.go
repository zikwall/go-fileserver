package actions

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (a *ActionProvider) PullFile(ctx *fiber.Ctx) error {
	filename, ok := ctx.Locals(a.FilenameKey).(string)

	if !ok {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": fmt.Errorf("Can't find file: %v", filename),
		})
	}

	return ctx.SendFile(fmt.Sprintf("%s/%s", a.RootFileDirectory, filename))
}

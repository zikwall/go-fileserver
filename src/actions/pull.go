package actions

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/lib"
)

func (a *ActionProvider) PullFile(ctx *fiber.Ctx) error {
	filename, ok := ctx.Locals(a.FilenameKey).(string)

	if !ok {
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: fmt.Sprintf("Can't find file: %v", filename),
		})
	}

	file := fmt.Sprintf("%s/%s", a.RootFileDirectory, filename)

	if !lib.Exists(file) {
		return ctx.Status(404).JSON(Response{
			Status:  false,
			Message: fmt.Sprintf("File not found: %s", filename),
		})
	}

	return ctx.SendFile(file)
}

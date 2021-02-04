package actions

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/lib"
)

type ActionProvider struct {
	FilenameKey       string
	FormFilesKey      string
	FormFileKey       string
	RootFileDirectory string
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func (a ActionProvider) GetFilepath(ctx *fiber.Ctx) (string, string, int, error) {
	filename, ok := ctx.Locals(a.FilenameKey).(string)

	if !ok {
		return "", "", 500, fmt.Errorf("Can't find file: %v", filename)
	}

	file := fmt.Sprintf("%s/%s", a.RootFileDirectory, filename)

	if !lib.Exists(file) {
		return "", "", 404, fmt.Errorf("File not found: %s", filename)
	}

	return file, filename, 200, nil
}

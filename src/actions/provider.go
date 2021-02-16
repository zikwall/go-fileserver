package actions

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/lib"
	"path"
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

	filepath := path.Join(a.RootFileDirectory, filename)

	if !lib.Exists(filepath) {
		return "", "", 404, fmt.Errorf("File not found: %s", filename)
	}

	return filepath, filename, 200, nil
}

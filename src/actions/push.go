package actions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/lib"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func (a *ActionProvider) PushFile(ctx *fiber.Ctx) error {
	filename, ok := ctx.Locals("filename").(string)

	if !ok {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": fmt.Errorf("Can't push file: %v", filename),
		})
	}

	if strings.Contains(ctx.Get("Content-Type"), "multipart/form-data") {
		// Try read one file, example -F "file=@.gitignore"
		//
		// example:
		// ```bash
		// $ curl -i -X POST -H "Content-Type: multipart/form-data" -F "file=@.gitignore" http://localhost:1337/.gitignore?token=123456
		// ```
		if file, err := ctx.FormFile(a.FormFileKey); err == nil {
			if err := ctx.SaveFile(file, fmt.Sprintf("%s/%s", a.RootFileDirectory, file.Filename)); err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"status":  false,
					"message": fmt.Errorf("Failed save file: %s with error: %s", file.Filename, err),
				})
			}

			return ctx.JSON(fiber.Map{
				"status":  true,
				"message": "Successfully upload file!",
			})
		}

		// Handle multiple files, example -F "files[]=@.gitignore"
		//
		// example:
		// ```bash
		// $ curl -i -X POST -H "Content-Type: multipart/form-data" -F "files[]=@.gitignore" -F "files[]=@README.md" http://localhost:1337/.gitignore?token=123456
		// ```
		form, err := ctx.MultipartForm()

		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": fmt.Errorf("Failed parse multipart/form-data, error: %s", err).Error(),
			})
		}

		files, ok := form.File[a.FormFilesKey]

		if !ok {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": errors.New("Failed parse multipart/form-data: not contains `files[]` field").Error(),
			})
		}

		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			if err := ctx.SaveFile(file, fmt.Sprintf("%s/%s", a.RootFileDirectory, file.Filename)); err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"status":  false,
					"message": fmt.Errorf("Failed save file: %s wit error: %s", file.Filename, err).Error(),
				})
			}
		}

		return ctx.JSON(fiber.Map{
			"status":  true,
			"message": "Successfully upload all files!",
		})
	}

	// Read full body bytes, example --data-binary format
	//
	// example:
	// ```bash
	// $ curl -i -X POST --data-binary @.gitignore http://localhost:1337/.gitignore?token=123456
	// ```

	tempFile, err := ioutil.TempFile(a.RootFileDirectory, "upload_*")

	defer func() {
		if lib.Exists(tempFile.Name()) {
			if err := os.Remove(tempFile.Name()); err != nil {
				lib.Warning(err)
			}
		}
	}()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	reader := bytes.NewReader(ctx.Body())
	_, err = io.Copy(tempFile, reader)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	// You must explicitly close it, otherwise there will be an error:
	// The process cannot access the file because it is being used by another process
	if err := tempFile.Close(); err != nil {
		lib.Warning(err)
	}

	if err := os.Rename(tempFile.Name(), fmt.Sprintf("%s/%s", a.RootFileDirectory, filename)); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  true,
		"message": "Successfully upload file!",
	})
}

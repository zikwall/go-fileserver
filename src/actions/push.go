package actions

import (
	"bytes"
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
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: fmt.Sprintf("Can't push file: %v", filename),
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
				return ctx.Status(500).JSON(Response{
					Status:  false,
					Message: fmt.Sprintf("Failed save `%s`, error: %s", file.Filename, err),
				})
			}

			return ctx.JSON(Response{
				Status:  true,
				Message: "Successfully upload file!",
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
			return ctx.Status(500).JSON(Response{
				Status:  false,
				Message: err.Error(),
			})
		}

		files, ok := form.File[a.FormFilesKey]

		if !ok {
			return ctx.Status(500).JSON(Response{
				Status:  false,
				Message: "Failed parse multipart/form-data: not contains `files[]` field",
			})
		}

		for _, file := range files {
			filepath := fmt.Sprintf("%s/%s", a.RootFileDirectory, file.Filename)

			if err := ctx.SaveFile(file, filepath); err != nil {
				return ctx.Status(500).JSON(Response{
					Status:  false,
					Message: fmt.Sprintf("Failed save `%s`, error: %s", file.Filename, err),
				})
			}
		}

		return ctx.JSON(Response{
			Status:  true,
			Message: "Successfully upload all files!",
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
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: err.Error(),
		})
	}

	reader := bytes.NewReader(ctx.Body())
	_, err = io.Copy(tempFile, reader)

	if err != nil {
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: err.Error(),
		})
	}

	// You must explicitly close it, otherwise there will be an error:
	// The process cannot access the file because it is being used by another process
	if err := tempFile.Close(); err != nil {
		lib.Warning(err)
	}

	filepath := fmt.Sprintf("%s/%s", a.RootFileDirectory, filename)

	if err := os.Rename(tempFile.Name(), filepath); err != nil {
		return ctx.Status(500).JSON(Response{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.JSON(Response{
		Status:  true,
		Message: "Successfully upload file!",
	})
}

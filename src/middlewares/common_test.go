package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/actions"
	"github.com/zikwall/go-fileserver/src/lib"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCommon(t *testing.T) {
	t.Run("it should be correct work", func(t *testing.T) {
		absolutePath, err := filepath.Abs("../../tmp")

		if err != nil {
			t.Fatal(err)
		}

		token, err := lib.GenerateToken()

		if err != nil {
			t.Fatal(err)
		}

		action := actions.ActionProvider{
			FilenameKey:       "filename",
			FormFilesKey:      "files[]",
			FormFileKey:       "file",
			RootFileDirectory: absolutePath,
		}

		app := fiber.New()
		app.Use(WithFilename())
		app.Use(WithProtection(token))

		app.Get("/:filename", action.PullFile)

		app.Put("/:filename?",
			WithPushable(),
			action.PushFile,
		)

		app.Post("/:filename?",
			WithPushable(),
			action.PushFile,
		)

		t.Run("push with query params", func(t *testing.T) {
			temp, err := ioutil.TempFile("./", "test_file_*.txt")

			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = temp.Close()
				_ = os.Remove(temp.Name())
			}()

			text := []byte("This is a Go!")

			if _, err = temp.Write(text); err != nil {
				t.Fatal(err)
			}

			t.Run("POST multiple files with query params", func(t *testing.T) {
				req, err := createMultipartRequest(fmt.Sprintf("/test_file.txt?token=%s", token), "files[]", temp.Name(), map[string]string{})

				if err != nil {
					t.Fatal(err)
				}

				// q := req.URL.Query()
				// q.Add("token", token)
				// req.URL.RawQuery = q.Encode()

				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET previous file with HEADER", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", temp.Name()), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET status 401 previous file with wrong HEADER", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", temp.Name()), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", "wrong_token_here"))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 401 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET status 404 previous file with wrong filename", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", "wrong_file_name.txt"), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 404 {
					t.Fatal("Failed check signature by response status code")
				}
			})
		})

		t.Run("push with form value", func(t *testing.T) {
			temp, err := ioutil.TempFile("./", "test_file_*.txt")

			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = temp.Close()
				_ = os.Remove(temp.Name())
			}()

			text := []byte("This is a Go!")

			if _, err = temp.Write(text); err != nil {
				t.Fatal(err)
			}

			t.Run("POST multiple files", func(t *testing.T) {
				req, err := createMultipartRequest("/test_file.txt", "files[]", temp.Name(), map[string]string{
					"token": token,
				})

				if err != nil {
					t.Fatal(err)
				}

				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET previous file with HEADER", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", temp.Name()), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})
		})

		t.Run("push with header", func(t *testing.T) {
			temp, err := ioutil.TempFile("./", "test_file_*.txt")

			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = temp.Close()
				_ = os.Remove(temp.Name())
			}()

			text := []byte("This is a Go!")

			if _, err = temp.Write(text); err != nil {
				t.Fatal(err)
			}

			t.Run("POST multiple files", func(t *testing.T) {
				req, err := createMultipartRequest("/test_file.txt", "files[]", temp.Name(), map[string]string{})

				if err != nil {
					t.Fatal(err)
				}

				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))

				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET previous file with HEADER", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", temp.Name()), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})
		})

		t.Run("push with header", func(t *testing.T) {
			temp, err := ioutil.TempFile("./", "test_file_*.txt")

			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				_ = temp.Close()
				_ = os.Remove(temp.Name())
			}()

			text := []byte("This is a Go!")

			if _, err = temp.Write(text); err != nil {
				t.Fatal(err)
			}

			t.Run("POST multiple files", func(t *testing.T) {
				req, err := createMultipartRequest("/test_file.txt", "files[]", temp.Name(), map[string]string{})

				if err != nil {
					t.Fatal(err)
				}

				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))

				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})

			t.Run("it should be successfully GET previous file with HEADER", func(t *testing.T) {
				req := httptest.NewRequest("GET", fmt.Sprintf("/%s", temp.Name()), nil)
				req.Header.Set(AuthHeader, fmt.Sprintf("Bearer %s", token))
				resp, err := app.Test(req)

				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != 200 {
					t.Fatal("Failed check signature by response status code")
				}
			})
		})
	})
}

func createMultipartRequest(uri string, paramName, path string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

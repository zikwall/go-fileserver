package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"github.com/zikwall/go-fileserver/src/actions"
	"github.com/zikwall/go-fileserver/src/lib"
	"github.com/zikwall/go-fileserver/src/middlewares"
	"log"
	"os"
	"path/filepath"
)

func main() {
	application := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bind-address",
				Value:   "0.0.0.0:1337",
				Usage:   "Run service in host",
				EnvVars: []string{"BIND_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "Token to protect requests (if empty is generated automatically)",
				EnvVars: []string{"TOKEN"},
			},
			&cli.StringFlag{
				Name:    "root-file-directory",
				Usage:   "",
				EnvVars: []string{"ROOT_FILE_DIRECTORY"},
				Value:   "./tmp",
			},
			&cli.BoolFlag{
				Name:    "enable-secure",
				Usage:   "Enabling/disabling token protection",
				EnvVars: []string{"ENABLE_SECURE"},
				Value:   true,
			},
		},
	}

	application.Action = func(c *cli.Context) error {
		app := fiber.New()
		app.Use(middlewares.WithFilename())

		if c.Bool("enable-secure") {
			token := c.String("token")

			if token == "" {
				generated, err := generateToken()

				if err != nil {
					return err
				}

				token = generated

				lib.Info(fmt.Sprintf("TOKEN: %s", generated))
			}

			app.Use(middlewares.WithProtection(token))
		}

		absolutePath, err := filepath.Abs(c.String("root-file-directory"))

		if err != nil {
			return err
		}

		action := actions.ActionProvider{
			FilenameKey:       "filename",
			FormFilesKey:      "files[]",
			FormFileKey:       "file",
			RootFileDirectory: absolutePath,
		}

		app.Get("/:filename", action.PullFile)

		app.Put("/:filename?",
			middlewares.WithPushable(),
			action.PushFile,
		)

		app.Post("/:filename?",
			middlewares.WithPushable(),
			action.PushFile,
		)

		go func() {
			if err := app.Listen(c.String("bind-address")); err != nil {
				log.Fatal(err)
			}
		}()

		congratulations()
		waitSystemNotify()

		return nil
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

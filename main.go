package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"github.com/zikwall/go-fileserver/src/actions"
	"github.com/zikwall/go-fileserver/src/constants"
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
				Value:   false,
			},
			&cli.IntFlag{
				Name:    "secure-type",
				Usage:   "Token=0, Basic auth=1 or JWT=2",
				EnvVars: []string{"SECURE_TYPE"},
				Value:   0,
			},
			&cli.StringSliceFlag{
				Name:    "users",
				Usage:   "Users, format username:password",
				EnvVars: []string{"USERS"},
			},
			&cli.BoolFlag{
				Name:    "enable-tsl",
				Usage:   "Enabling/disabling SSL protection",
				EnvVars: []string{"ENABLE_TSL"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "tsl-cert-file",
				Usage:   "",
				EnvVars: []string{"TSL_CERT_FILE"},
				Value:   "./tmp/ssl.cert",
			},
			&cli.StringFlag{
				Name:    "tsl-key-file",
				Usage:   "",
				EnvVars: []string{"TSL_KEY_FILE"},
				Value:   "./tmp/ssl.key",
			},
		},
	}

	application.Action = func(c *cli.Context) error {
		app := fiber.New()
		app.Use(middlewares.WithFilename())

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

		if c.Bool("enable-secure") {
			switch c.Int("secure-type") {
			case constants.TypeToken:
				app.Use(middlewares.WithProtection(c.String("token")))
			case constants.TypeBasic:
				app.Use(middlewares.WithBasicAuth(c.StringSlice("users")...))
			case constants.TypeJWT:
				fallthrough
			default:
				log.Fatalf("Unsupported secure type: %d", c.Int("secure-type"))
			}
		}

		app.Delete("/:filename", action.DeleteFile)

		app.Put("/:filename?",
			middlewares.WithPushable(),
			action.PushFile,
		)

		app.Post("/:filename?",
			middlewares.WithPushable(),
			action.PushFile,
		)

		go func() {
			if c.Bool("enable-tsl") {
				err := app.ListenTLS(
					c.String("bind-address"),
					c.String("tsl-cert-file"),
					c.String("tsl-key-file"),
				)

				if err != nil {
					log.Fatal(err)
				}
			} else {
				if err := app.Listen(c.String("bind-address")); err != nil {
					log.Fatal(err)
				}
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

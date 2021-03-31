package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"strings"
)

func WithBasicAuth(u ...string) fiber.Handler {
	users := make(map[string]string, len(u))

	for _, user := range u {
		userPars := strings.Split(user, ":")

		if len(userPars) != 2 {
			continue
		}

		users[userPars[0]] = userPars[1]
	}

	return basicauth.New(basicauth.Config{
		Users: users,
	})
}

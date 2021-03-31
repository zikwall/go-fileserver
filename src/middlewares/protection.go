package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zikwall/go-fileserver/src/lib"
	"strings"
)

const (
	AuthHeader     = "Authorization"
	AuthTokenType  = "Bearer"
	QueryTokenName = "token"
)

var (
	ErrorMissingToken  = errors.New("Missing token")
	ErrorTokenMismatch = errors.New("Mismatch token")
)

func authHeaderToken(header string) (string, bool) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], AuthTokenType) {
		return "", false
	}

	return parts[1], true
}

func WithProtection(secureToken string) fiber.Handler {
	if secureToken == "" {
		generated, _ := lib.GenerateToken()
		secureToken = generated

		lib.Info(fmt.Sprintf("Generate token: %s", generated))
	}

	return func(ctx *fiber.Ctx) error {
		token, ok := authHeaderToken(ctx.Get(AuthHeader))

		if !ok || token == "" {
			token = ctx.Query(QueryTokenName, "")

			if token == "" {
				token = ctx.FormValue(QueryTokenName, "")
			}
		}

		if token == "" {
			return ctx.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": ErrorMissingToken.Error(),
			})
		}

		if !strings.EqualFold(token, secureToken) {
			return ctx.Status(401).JSON(fiber.Map{
				"status":  false,
				"message": ErrorTokenMismatch.Error(),
			})
		}

		return ctx.Next()
	}
}

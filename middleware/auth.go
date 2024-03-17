package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ramadhan1445sprint/sprint_segokuning/crypto"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
)

func Auth(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	if auth == "" {
		return customErr.NewUnauthorizedError("token not found")
	}

	splitted := strings.Split(auth, " ")

	if splitted[0] != "Bearer" {
		return customErr.NewUnauthorizedError("invalid token")
	}

	payload, err := crypto.VerifyToken(splitted[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return customErr.NewUnauthorizedError("token expired")
		}
		return customErr.NewUnauthorizedError(err.Error())
	}

	ctx.Locals("user_id", payload.Id)
	ctx.Locals("cred_type", payload.Type)
	ctx.Locals("cred_value", payload.Value)
	ctx.Locals("name", payload.Name)
	return ctx.Next()
}

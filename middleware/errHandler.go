package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := ctx.Response().StatusCode()
	if customError, ok := err.(customErr.CustomError); ok {
		code = customError.Status()
		return ctx.Status(code).JSON(fiber.Map{
			"message": customError.Error(),
		})
	} else if code < 400 {
		code = fiber.StatusInternalServerError
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": err.Error(),
	})
}

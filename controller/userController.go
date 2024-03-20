package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type UserController struct {
	svc svc.UserSvc
}

func NewUserController(svc svc.UserSvc) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) UpdateAccountUser(ctx *fiber.Ctx) error {
	var user entity.UpdateAccountPayload
	userId := ctx.Locals("user_id").(string)

	fmt.Println(userId)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "body parsing error"})
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	err := c.svc.UpdateAccountUser(user, userId)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "account updated successfully",
	})
}

func (c *UserController) UpdateLinkEmailAccount(ctx *fiber.Ctx) error {
	var user entity.PostLinkEmailPayload
	userId := ctx.Locals("user_id").(string)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "body parsing error"})
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	err := c.svc.UpdateLinkEmailAccount(user.Email, userId)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "email linked successfully",
	})
}

func (c *UserController) UpdateLinkPhoneAccount(ctx *fiber.Ctx) error {
	var user entity.PostLinkPhonePayload
	userId := ctx.Locals("user_id").(string)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "body parsing error"})
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	err := c.svc.UpdateLinkPhoneAccount(user.Phone, userId)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "phone linked successfully",
	})
}

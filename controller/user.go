package controller

import (
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

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var newUser entity.RegistrationPayload
	if err := ctx.BodyParser(&newUser); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	accessToken, err := c.svc.Register(newUser)
	if err != nil {
		return err
	}

	respData := fiber.Map{
		"name":        newUser.Name,
		"accessToken": accessToken,
	}

	if newUser.CredentialType == entity.Email {
		respData["email"] = newUser.CredentialValue
	} else if newUser.CredentialType == entity.Phone {
		respData["phone"] = newUser.CredentialValue
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    respData,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var creds entity.Credential
	if err := ctx.BodyParser(&creds); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	user, accessToken, err := c.svc.Login(creds)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"accessToken": accessToken,
			"name":        user.Name,
			"phone":       user.Phone,
			"email":       user.Email,
		},
	})
}

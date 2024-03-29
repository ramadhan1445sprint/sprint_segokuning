package controller

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type PostController struct {
	svc      svc.PostSvc
	validate *validator.Validate
}

func NewPostController(svc svc.PostSvc, validate *validator.Validate) *PostController {
	return &PostController{svc: svc, validate: validate}
}

func (c *PostController) CreatePost(ctx *fiber.Ctx) error {
	var postReq entity.Post
	userId := ctx.Locals("user_id").(string)
	
	if err := ctx.BodyParser(&postReq); err != nil {
		custErr := customErr.NewBadRequestError(err.Error())
		return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
	}
	
	if err := c.validate.Struct(postReq); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			custErr := customErr.NewBadRequestError(e.Error())
			return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
		}
	}

	postReq.UserID = userId

	if err := c.svc.CreatePost(postReq); err != nil {
		return ctx.Status(err.StatusCode).JSON(fiber.Map{"message": err.Message})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
}

func (c *PostController) GetPost(ctx *fiber.Ctx) error {
	var filterReq entity.PostFilter

	limit := ctx.Query("limit")
	if limit == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid param",
		})
	}
	offset := ctx.Query("offset")
	if offset == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid param",
		})
	}

	if err := ctx.QueryParser(&filterReq); err != nil {
		custErr := customErr.NewBadRequestError(err.Error())
		return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
	}

	if err := c.validate.Struct(filterReq); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			custErr := customErr.NewBadRequestError(e.Error())
			return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
		}
	}

	log.Printf("%+v\n", filterReq)

	if filterReq.Limit == 0 {
		filterReq.Limit = 5
	}

	resp, err := c.svc.GetPost(filterReq)

	if err != nil {
		return ctx.Status(err.StatusCode).JSON(fiber.Map{"message": err.Message})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

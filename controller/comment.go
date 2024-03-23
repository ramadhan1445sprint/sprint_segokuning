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

type CommentController struct {
	svc svc.CommentSvc
	validate *validator.Validate
}

func NewCommentController(svc svc.CommentSvc, validate *validator.Validate) *CommentController {
	return &CommentController{svc: svc, validate: validate}
}

func (c *CommentController) CreateComment(ctx *fiber.Ctx) error {
	var commentReq entity.Comment
	userId := ctx.Locals("user_id").(string)

	if err := ctx.BodyParser(&commentReq); err != nil {
		custErr := customErr.NewBadRequestError(err.Error())
		return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
	}

	if err := c.validate.Struct(commentReq); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			custErr := customErr.NewBadRequestError(e.Error())
			return ctx.Status(custErr.StatusCode).JSON(fiber.Map{"message": custErr.Message})
		}
	}
	commentReq.UserID = userId
	log.Printf("%+v\n", commentReq)

	if err := c.svc.CreateComment(commentReq); err != nil {
		return ctx.Status(err.StatusCode).JSON(fiber.Map{"message": err.Message})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "success"})
}

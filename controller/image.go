package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type ImageController struct {
	svc svc.ImageSvc
}

func NewImageController(svc svc.ImageSvc) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (c *ImageController) UploadImage(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return customErr.NewInternalServerError("failed to retrieve file")
	}

	url, err := c.svc.UploadImage(fileHeader)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "File uploaded sucessfully",
		"data": fiber.Map{
			"imageUrl": url,
		},
	})

}

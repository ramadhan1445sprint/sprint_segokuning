package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type Controller struct {
	svc svc.Svc
}

func NewController(svc svc.Svc) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (c *Controller) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.SendString("Healthy")
}

func (c *Controller) AuthCheck(ctx *fiber.Ctx) error {
	return ctx.SendString("Authorized")
}

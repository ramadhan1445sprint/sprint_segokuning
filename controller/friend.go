package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type FriendController struct {
	svc svc.FriendSvc
}

func NewFriendController(svc svc.FriendSvc) FriendController {
	return FriendController{svc}
}

func (c *FriendController) AddFriend(ctx *fiber.Ctx) error {
	var payload entity.AddDeleteFriendPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	userId := ctx.Locals("user_id").(string)

	err := c.svc.AddFriend(userId, payload.FriendId)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Add Friend Success!",
	})
}

func (c *FriendController) DeleteFriend(ctx *fiber.Ctx) error {
	var payload *entity.AddDeleteFriendPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	userId := ctx.Locals("user_id").(string)

	err := c.svc.DeleteFriend(userId, payload.FriendId)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Delete Friend Success!",
	})
}

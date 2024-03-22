package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

type FriendController struct {
	svc svc.FriendSvc
}

func NewFriendController(svc svc.FriendSvc) *FriendController {
	return &FriendController{svc: svc}
}

func (c *FriendController) GetListFriends(ctx *fiber.Ctx) error {
	var param entity.ListFriendPayload
	ctx.QueryParser(&param)
	fmt.Printf("%v", param.OnlyFriend)

	param = entity.NewListFriend(param.OrderBy, param.OnlyFriend, param.SortBy, param.Search, param.Limit, param.Offset)

	if param.OnlyFriend {
		param.UserId = ctx.Locals("user_id").(string)
	}

	friends, meta, err := c.svc.GetListFriends(&param)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "",
		"data":    friends,
		"meta":    meta,
	})
}

package controller

import (
	"strconv"
	"strings"

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

func (c *FriendController) GetListFriends(ctx *fiber.Ctx) error {
	var err error

	params := ctx.Queries()

	var limit, offset int
	var onlyFriend bool
	var sortBy, orderBy, search string
	qs := string(ctx.Request().URI().QueryString())

	search = params["search"]

	if params["limit"] == "" {
		if strings.Contains(qs, "limit") {
			return customErr.NewBadRequestError("invalid query limit args")
		} else {
			params["limit"] = "5"
		}
	}

	limit, err = strconv.Atoi(params["limit"])
	if err != nil {
		return customErr.NewBadRequestError("invalid parsing query limit args to int")
	}
	if limit < 0 {
		return customErr.NewBadRequestError("limit cannot be negative")
	}

	if params["offset"] == "" {
		if strings.Contains(qs, "offset") {
			return customErr.NewBadRequestError("invalid query offset args")
		} else {
			params["offset"] = "0"
		}
	}

	offset, err = strconv.Atoi(params["offset"])
	if err != nil {
		return customErr.NewBadRequestError("invalid parsing query offset args to int")
	}
	if offset < 0 {
		return customErr.NewBadRequestError("limit cannot be negative")
	}

	sortBy = params["sortBy"]

	if sortBy == "" {
		if strings.Contains(qs, "sortBy") {
			return customErr.NewBadRequestError("invalid query sortBy args")
		} else {
			sortBy = "createdAt"
		}
	} else if sortBy != "createdAt" && sortBy != "friendCount" {
		return customErr.NewBadRequestError("invalid query sortBy args")
	}

	orderBy = params["orderBy"]

	if orderBy == "" {
		if strings.Contains(qs, "orderBy") {
			return customErr.NewBadRequestError("invalid query orderBy args")
		} else {
			orderBy = "desc"
		}
	} else if orderBy != "desc" && orderBy != "asc" {
		return customErr.NewBadRequestError("invalid query orderBy args")
	}

	if params["onlyFriend"] == "" {
		if strings.Contains(qs, "onlyFriend") {
			return customErr.NewBadRequestError("invalid query onlyFriend args")
		} else {
			params["onlyFriend"] = "false"
		}
	}

	onlyFriend, err = strconv.ParseBool(params["onlyFriend"])
	if err != nil {
		return customErr.NewBadRequestError("invalid parsing query onlyFriend args to int")
	}

	param := &entity.ListFriendPayload{
		Limit:      limit,
		Offset:     offset,
		Search:     search,
		OnlyFriend: onlyFriend,
		SortBy:     entity.FriendSortBy(sortBy),
		OrderBy:    orderBy,
	}

	if param.OnlyFriend {
		param.UserId = ctx.Locals("user_id").(string)
	}

	friends, meta, err := c.svc.GetListFriends(param)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "",
		"data":    friends,
		"meta":    meta,
	})
}

package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type sortBy string
type FriendSortBy sortBy

var (
	SortByCreatedAt   FriendSortBy = "createdAt"
	SortByFriendCount FriendSortBy = "friendCount"
)

var friendSortBys []interface{} = []interface{}{SortByCreatedAt, SortByFriendCount}

type ListFriendPayload struct {
	OnlyFriend bool `binding:"omitempty" query:"onlyFriend"`
	UserId     string
	Search     string       `schema:"search" binding:"omitempty"`
	Limit      int          `schema:"limit" binding:"omitempty"`
	Offset     int          `schema:"offset" binding:"omitempty"`
	SortBy     FriendSortBy `schema:"sortby" binding:"omitempty"`
	OrderBy    string       `schema:"orderby" binding:"omitempty"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func (r ListFriendPayload) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.When(r.OnlyFriend, validation.Required.Error("userId is required"))),
		validation.Field(&r.OrderBy, validation.In("asc", "desc")),
		validation.Field(&r.SortBy, validation.In(friendSortBys...)),
	)

	return err
}

func NewListFriend(orderBy string, onlyFriend bool, sortBy FriendSortBy, search string, limit int, offset int) ListFriendPayload {
	if sortBy == "" || sortBy == "createdAt" {
		sortBy = SortByCreatedAt
	} else {
		sortBy = SortByFriendCount
	}

	if orderBy == "" {
		orderBy = "desc"
	}

	if limit == 0 {
		limit = 5
	}

	return ListFriendPayload{
		OnlyFriend: onlyFriend,
		UserId:     "",
		Search:     search,
		Limit:      limit,
		Offset:     offset,
		SortBy:     sortBy,
		OrderBy:    orderBy,
	}
}

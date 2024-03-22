package svc

import (
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type FriendSvc interface {
	GetListFriends(*entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error)
}

type friendSvc struct {
	repo repo.FriendRepo
}

func NewFriendService(repo repo.FriendRepo) FriendSvc {
	return &friendSvc{repo}
}

func (r *friendSvc) GetListFriends(param *entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error) {
	// validate list
	var listUser []entity.UserList

	if err := param.Validate(); err != nil {
		return listUser, nil, customErr.NewBadRequestError(err.Error())
	}

	listUser, meta, err := r.repo.GetListFriends(param)
	if err != nil {
		return listUser, nil, err
	}

	return listUser, meta, err
}

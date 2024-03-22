package svc

import (
	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type FriendSvc interface {
	AddFriend(string, string) error
	DeleteFriend(string, string) error
	GetListFriends(*entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error)
}

type friendSvc struct {
	userRepo   repo.UserRepo
	friendRepo repo.FriendRepo
}

func NewFriendSvc(userRepo repo.UserRepo, friendRepo repo.FriendRepo) FriendSvc {
	return &friendSvc{userRepo, friendRepo}
}

func (s *friendSvc) AddFriend(userId, friendId string) error {
	if userId == friendId {
		return customErr.NewBadRequestError("cannot add yourself as friend")
	}

	// check id is valid
	err := uuid.Validate(friendId)
	if err != nil {
		return customErr.NewBadRequestError("invalid friend id format")
	}

	// check friend is exist
	_, err = s.userRepo.GetUserById(friendId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return customErr.NewNotFoundError("friend not found")
		}
		return err
	}

	// check friend is not yet user friend
	friendPair, err := s.friendRepo.GetFriendPair(userId, friendId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}

	if friendPair != nil {
		return customErr.NewBadRequestError("friend already added")
	}

	err = s.friendRepo.AddFriend(userId, friendId)
	if err != nil {
		return err
	}

	return nil
}

func (s *friendSvc) DeleteFriend(userId, friendId string) error {
	if userId == friendId {
		return customErr.NewBadRequestError("cannot delete yourself from friend")
	}

	// check id is valid
	err := uuid.Validate(friendId)
	if err != nil {
		return customErr.NewBadRequestError("invalid friend id format")
	}

	// check friend is exist
	_, err = s.userRepo.GetUserById(friendId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return customErr.NewNotFoundError("friend not found")
		}
		return err
	}

	// check friend is user friend
	_, err = s.friendRepo.GetFriendPair(userId, friendId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return customErr.NewBadRequestError("friend not yet added")
		}
		return err
	}

	err = s.friendRepo.DeleteFriend(userId, friendId)
	if err != nil {
		return err
	}

	return nil
}

func (r *friendSvc) GetListFriends(param *entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error) {
	// validate list
	var listUser []entity.UserList

	if err := param.Validate(); err != nil {
		return listUser, nil, customErr.NewBadRequestError(err.Error())
	}

	listUser, meta, err := r.friendRepo.GetListFriends(param)
	if err != nil {
		return listUser, nil, err
	}

	return listUser, meta, err
}

package svc

import (
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type UserSvc interface {
	UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error
	UpdateLinkEmailAccount(email string, userId string) error
	UpdateLinkPhoneAccount(phone string, userId string) error
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error {
	err := s.repo.UpdateAccountUser(user, userId)
	if err != nil {
		if err.Error() == "user not found" {
			return customErr.NewNotFoundError(err.Error())
		} else {
			return customErr.NewInternalServerError(err.Error())
		}
	}

	return nil
}

func (s *userSvc) UpdateLinkEmailAccount(email string, userId string) error {
	// check duplicate email
	if err := s.repo.GetUserId(email, "email"); err != nil {
		return err
	}

	// check user already linked email
	if err := s.repo.GetUser(userId); err != nil {
		return err
	}

	err := s.repo.UpdateLinkAccount(email, userId, "email")
	if err != nil {
		if err.Error() == "user not found" {
			return customErr.NewNotFoundError(err.Error())
		} else {
			return customErr.NewInternalServerError(err.Error())
		}
	}

	return nil
}

func (s *userSvc) UpdateLinkPhoneAccount(phone string, userId string) error {
	// check duplicate phone
	if err := s.repo.GetUserId(phone, "phone"); err != nil {
		return err
	}

	// check user already linked phone
	if err := s.repo.GetUser(userId); err != nil {
		return err
	}

	err := s.repo.UpdateLinkAccount(phone, userId, "phone")
	if err != nil {
		if err.Error() == "user not found" {
			return customErr.NewNotFoundError(err.Error())
		} else {
			return customErr.NewInternalServerError(err.Error())
		}
	}

	return nil
}

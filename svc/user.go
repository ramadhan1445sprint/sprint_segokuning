package svc

import (
	"fmt"

	"github.com/ramadhan1445sprint/sprint_segokuning/crypto"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type UserSvc interface {
	Login(creds entity.Credential) (*entity.User, error)
	Register(user entity.RegistrationPayload) (string, error)
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) Login(creds entity.Credential) (*entity.User, error) {
	return &entity.User{}, nil
}

func (s *userSvc) Register(payload entity.RegistrationPayload) (string, error) {
	// check if credential type not phone or email
	if payload.CredentialType != entity.Email && payload.CredentialType != entity.Phone {
		return "", customErr.NewBadRequestError("credential type must be email or phone")
	}

	// check if credential value is exists
	existingUser, err := s.repo.GetUser(payload.CredentialValue, string(payload.CredentialType))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return "", err
		}
	}

	if existingUser != nil {
		errMsg := fmt.Sprintf("user with %s %s already exists", payload.CredentialType, payload.CredentialValue)
		return "", customErr.NewConflictError(errMsg)
	}

	user := entity.NewUser(payload.CredentialType, payload.CredentialValue, payload.Name, payload.Password)

	err = user.Validate(payload.CredentialType)
	if err != nil {
		return "", customErr.NewBadRequestError(err.Error())
	}

	hashedPassword, err := crypto.GenerateHashedPassword(payload.Password)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(&payload, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err := crypto.GenerateToken(id, payload.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

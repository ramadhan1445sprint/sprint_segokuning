package svc

import (
	"fmt"

	"github.com/ramadhan1445sprint/sprint_segokuning/crypto"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type UserSvc interface {
	Login(creds entity.Credential) (*entity.User, string, error)
	Register(user entity.RegistrationPayload) (string, error)
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

func (s *userSvc) Login(creds entity.Credential) (*entity.User, string, error) {
	err := creds.Validate()
	if err != nil {
		return nil, "", customErr.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(creds.CredentialValue, creds.CredentialType)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, "", customErr.NewNotFoundError("user not found")
		}
		return nil, "", err
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return nil, "", customErr.NewBadRequestError("wrong password!")
	}

	token, err := crypto.GenerateToken(user.Id, user.Name)
	if err != nil {
		return nil, "", customErr.NewBadRequestError(err.Error())
	}

	return user, token, nil
}

func (s *userSvc) Register(payload entity.RegistrationPayload) (string, error) {
	if payload.CredentialType != entity.Email && payload.CredentialType != entity.Phone {
		return "", customErr.NewBadRequestError("credential type must be email or phone")
	}

	existingUser, err := s.repo.GetUser(payload.CredentialValue, payload.CredentialType)
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

func (s *userSvc) UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error {
	if err := user.Validate(); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	if err := s.repo.UpdateAccountUser(user, userId); err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *userSvc) UpdateLinkEmailAccount(email string, userId string) error {
	// validate email
	if err := entity.ValidateEmail(email); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	// check duplicate email
	existingUser, err := s.repo.GetUser(email, "email")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}

	if existingUser != nil {
		return customErr.NewConflictError("email already exists")
	}

	// check user already linked email
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Email != "" {
		return customErr.NewBadRequestError("email already linked")
	}

	if err := s.repo.UpdateLinkAccount(email, userId, "email"); err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *userSvc) UpdateLinkPhoneAccount(phone string, userId string) error {
	// validate phone
	if err := entity.ValidatePhone(phone); err != nil {
		return customErr.NewBadRequestError(err.Error())
	}

	// check duplicate phone
	existingUser, err := s.repo.GetUser(phone, "phone")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}

	if existingUser != nil {
		return customErr.NewConflictError("phone already exist")
	}

	// check user already linked phone
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Phone != "" {
		return customErr.NewBadRequestError("phone already linked")
	}

	if err := s.repo.UpdateLinkAccount(phone, userId, "phone"); err != nil {
		return customErr.NewInternalServerError(err.Error())
	}

	return nil
}

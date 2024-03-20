package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type UserRepo interface {
	GetUser(string, entity.CredType) (*entity.User, error)
	CreateUser(*entity.RegistrationPayload, string) (string, error)
	UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error
	GetUserById(userId string) (*entity.User, error)
	UpdateLinkAccount(credential string, userId string, credentialType string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) GetUser(credValue string, credType entity.CredType) (*entity.User, error) {
	var user entity.User

	statement := "SELECT id, name, email, phone, password, image_url, created_at FROM users"

	if credType == entity.Email {
		statement += " WHERE email = $1;"
	} else if credType == entity.Phone {
		statement += " WHERE phone = $1;"
	} else {
		return nil, customErr.NewBadRequestError("credential type must be email or phone")
	}

	err := r.db.Get(&user, statement, credValue)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) CreateUser(user *entity.RegistrationPayload, hashPassword string) (string, error) {
	var id string

	statement := fmt.Sprintf("INSERT INTO users (name, %s, password) VALUES ($1, $2, $3) RETURNING id", user.CredentialType)

	row := r.db.QueryRowx(
		statement,
		user.Name, user.CredentialValue, hashPassword,
	)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error {
	query := "UPDATE users SET name = $1, image_url = $2 WHERE id = $3"

	_, err := r.db.Exec(query, user.Name, user.ImageUrl, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUserById(userId string) (*entity.User, error) {
	var user entity.User

	query := "SELECT id, name, email, phone, password, image_url, created_at FROM users WHERE id = $1"

	err := r.db.Get(&user, query, userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateLinkAccount(credential string, userId string, credentialType string) error {
	var query string

	if credentialType == "phone" {
		query = "UPDATE users SET phone = $1 WHERE id = $2"
	} else if credentialType == "email" {
		query = "UPDATE users SET email = $1 WHERE id = $2"
	}

	_, err := r.db.Exec(query, credential, userId)
	if err != nil {
		return err
	}

	return nil
}

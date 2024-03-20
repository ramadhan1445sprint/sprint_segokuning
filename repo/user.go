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

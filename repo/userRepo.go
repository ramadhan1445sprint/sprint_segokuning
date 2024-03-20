package repo

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type UserRepo interface {
	UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error
	GetUserId(credential string, credentialType string) error
	GetUser(userId string) error
	UpdateLinkAccount(credential string, userId string, credentialType string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) UpdateAccountUser(user entity.UpdateAccountPayload, userId string) error {
	query := "UPDATE users SET name = $1, image_url = $2 WHERE id = $3"

	res, err := r.db.Exec(query, user.Name, user.ImageUrl, userId)
	if err != nil {
		return err
	}

	rowsEffected, _ := res.RowsAffected()

	if rowsEffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepo) GetUserId(credential string, credentialType string) error {
	var userId string
	var query string

	if credentialType == "phone" {
		query = "SELECT id FROM users WHERE phone = $1"
	} else if credentialType == "email" {
		query = "SELECT id FROM users WHERE email = $1"
	}

	err := r.db.QueryRow(query, credential).Scan(&userId)
	if userId != "" {
		return customErr.NewConflictError("credential already existed")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUser(userId string) error {
	query := "SELECT email, phone FROM users WHERE id = $1"

	row := r.db.QueryRow(query, userId)
	user := &entity.LinkAccountDetail{}

	err := row.Scan(&user.Email, &user.Phone)
	if user.Email != "" {
		return customErr.NewBadRequestError("user already linked email")
	} else if user.Phone != "" {
		return customErr.NewBadRequestError("user already linked phone")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdateLinkAccount(credential string, userId string, credentialType string) error {
	var query string

	if credentialType == "phone" {
		query = "UPDATE users SET phone = $1 WHERE id = $2"
	} else if credentialType == "email" {
		query = "UPDATE users SET email = $1 WHERE id = $2"
	}

	res, err := r.db.Exec(query, credential, userId)
	if err != nil {
		return err
	}

	rowsEffected, _ := res.RowsAffected()

	if rowsEffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

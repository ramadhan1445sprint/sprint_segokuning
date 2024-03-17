package repo

import "github.com/jmoiron/sqlx"

type Repo interface {}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}
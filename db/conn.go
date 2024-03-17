package db

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/config"
)

func NewDatabase() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_NAME"),
	)

	env := config.GetString("ENV")

	if env == "production" {
		dsn += " sslmode=verify-full rootcert=ap-southeast-1-bundle.pem"
	} else {
		dsn += " sslmode=disable"
	}

	db, err := sqlx.Connect("pgx", dsn)

	return db, err
}
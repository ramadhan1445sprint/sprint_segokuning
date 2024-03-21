package db

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/config"
)

func NewDatabase() (*sqlx.DB, error) {
	sslmode := config.GetString("DB_PARAMS")
	sslmode = strings.TrimPrefix(sslmode, "&")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_NAME"),
		sslmode,
	)

	if sslmode == "sslmode=verify-full" {
		dsn += " rootcert=ap-southeast-1-bundle.pem"
	}

	db, err := sqlx.Connect("pgx", dsn)

	return db, err
}
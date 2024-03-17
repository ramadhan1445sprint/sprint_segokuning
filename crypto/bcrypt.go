package crypto

import (
	"strconv"

	"github.com/ramadhan1445sprint/sprint_segokuning/config"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password string) (string, error) {
	costStr := config.GetString("BCRYPT_SALT")
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
package crypto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ramadhan1445sprint/sprint_segokuning/config"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

func GenerateToken(id, name string) (string, error) {
	secret := config.GetString("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JWTClaims{
		Id:              id,
		Name:            name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func VerifyToken(token string) (*entity.JWTPayload, error) {
	secret := config.GetString("JWT_SECRET")

	claims := &entity.JWTClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, err
	}
	
	payload := &entity.JWTPayload{
		Id:              claims.Id,
		Name:            claims.Name,
	}

	return payload, nil
}

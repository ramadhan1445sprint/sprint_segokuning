package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type Credential struct {
	Type     string `json:"credentialType"`
	Value    string `json:"credentialValue"`
	Password string `json:"password"`
}

type JWTPayload struct {
	Id    string
	Type  string
	Value string
	Name  string
}

type JWTClaims struct {
	Id    string
	Type  string
	Value string
	Name  string
	jwt.RegisteredClaims
}

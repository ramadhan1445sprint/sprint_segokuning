package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type Credential struct {
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
	Password        string `json:"password"`
}

type JWTPayload struct {
	Id              string
	CredentialType  string
	CredentialValue string
	Name            string
}

type JWTClaims struct {
	Id              string
	CredentialType  string
	CredentialValue string
	Name            string
	jwt.RegisteredClaims
}

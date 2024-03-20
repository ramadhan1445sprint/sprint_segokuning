package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type RegistrationPayload struct {
	CredentialType  CredType `json:"credentialType"`
	CredentialValue string   `json:"credentialValue"`
	Name            string   `json:"name"`
	Password        string   `json:"password"`
}

type Credential struct {
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
	Password        string `json:"password"`
}

type JWTPayload struct {
	Id   string
	Name string
}

type JWTClaims struct {
	Id   string
	Name string
	jwt.RegisteredClaims
}

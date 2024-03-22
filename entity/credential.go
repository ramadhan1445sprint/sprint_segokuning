package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
)

type RegistrationPayload struct {
	CredentialType  CredType `json:"credentialType"`
	CredentialValue string   `json:"credentialValue"`
	Name            string   `json:"name"`
	Password        string   `json:"password"`
}

type Credential struct {
	CredentialType  CredType `json:"credentialType"`
	CredentialValue string   `json:"credentialValue"`
	Password        string   `json:"password"`
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

func (c *Credential) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.CredentialType,
			validation.Required.Error("credential type is required"),
			validation.In(Phone, Email).Error("credential type must be phone or email"),
		),
		validation.Field(&c.CredentialValue,
			validation.When(c.CredentialType == Phone,
				validation.Required.Error("phone number is required"),
				validation.Length(7, 13).Error("phone number must between 7 and 13 digits including the country code"),
				validation.By(validatePhoneNumberFormat),
			).Else(
				validation.Required.Error("email is required"),
				validation.By(validateEmailFormat),
			),
		),
		validation.Field(&c.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 15).Error("password must be between 5 and 15 characters"),
		),
	)

	return err
}

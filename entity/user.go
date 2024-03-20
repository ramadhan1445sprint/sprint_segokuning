package entity

import (
	"errors"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CredType string

type User struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Phone       string `db:"phone"`
	Password    string `db:"password"`
	ImageUrl    string `db:"image_url"`
	FriendCount string `db:"friend_count"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

const (
	Phone CredType = "phone"
	Email CredType = "email"
)

func NewUser(credType CredType, credValue, name, password string) *User {
	u := &User{
		Name:     name,
		Password: password,
	}

	if credType == Phone {
		u.Phone = credValue
	} else {
		u.Email = credValue
	}

	return u
}

func (u *User) Validate(credentialType CredType) error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(5, 50).Error("name must be between 5 and 50 characters"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 15).Error("password must be between 5 and 15 characters"),
		),
		validation.Field(&u.Phone,
			validation.When(credentialType == Phone,
				validation.Required.Error("phone number is required"),
				validation.Length(7, 13).Error("phone number must between 7 and 13 digits including the country code"),
				validation.By(validatePhoneNumberFormat),
			),
		),
		validation.Field(&u.Email,
			validation.When(credentialType == Email,
				validation.Required.Error("email is required"),
				is.Email.Error("invalid email format"),
			),
		),
	)

	return err
}

func ValidateName(name string) error {
	err := validation.Validate(name,
		validation.Required.Error("name is required"),
		validation.Length(5, 50).Error("name must be between 5 and 50 characters"),
	)

	return err
}

func validatePhoneNumberFormat(value any) error {
	phoneNumber, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	pattern := `^\+(9[976]\d|8[987530]\d|6[987]\d|5[90]\d|42\d|3[875]\d|2[98654321]\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|4[987654310]|3[9643210]|2[70]|7|1)\d{1,}$`
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(phoneNumber) {
		return errors.New("invalid phone number format")
	}

	return nil
}

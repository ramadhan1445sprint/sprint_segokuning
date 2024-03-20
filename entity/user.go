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

type UpdateAccountPayload struct {
	ImageUrl string `db:"image_url" json:"imageUrl" validate:"required,url"`
	Name     string `db:"name" json:"name" validate:"required,max=50,min=5"`
}

type PostLinkEmailPayload struct {
	Email string `db:"email" json:"email,omitempty" validate:"required,email"`
}

type PostLinkPhonePayload struct {
	Phone string `db:"phone" json:"phone,omitempty" validate:"required,max=7,min=13,e164"`
}

type LinkAccountDetail struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

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

func validateEmailFormat(value any) error {
	email, ok := value.(string)
	if !ok {
		return errors.New("parse error")
	}

	pattern := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	rgx := regexp.MustCompile(pattern)
	if !rgx.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

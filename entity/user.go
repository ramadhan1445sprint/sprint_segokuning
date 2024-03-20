package entity

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

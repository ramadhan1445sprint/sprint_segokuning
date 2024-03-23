package entity

type Comment struct {
	ID   string   `json:"id"`
	UserID string  `json:"userId"`
	PostID string  `json:"postId" validate:"required"`
	Comment  string `json:"comment" validate:"required,min=2,max=500"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
package entity

type Comment struct {
	ID   string   `json:"id"`
	UserID string  `json:"userId"`
	PostID string  `json:"postId"`
	Comment  string `json:"comment"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
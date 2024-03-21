package entity

type Post struct {
	ID         string   `json:"postId"`
	UserID     string   `json:"user_id"`
	PostInHtml   string  `json:"postInHtml" db:"post_in_html"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt" db:"created_at"`
	UpdatedAt  string   `json:"updatedAt,omitempty"`
}

type PostData struct {
    PostId     string    `db:"post_id"`
    PostInHtml string    `db:"post_in_html"`
    Tags       []string  `db:"tags"`
    CreatedAt  string `db:"post_created_at"`
    Creator    PostUser `db:"creator"`
    Comments   []PostComment `db:"commentx"`
}

type PostUser struct {
    UserId      string `db:"user_id"`
    Name        string `db:"name"`
    ImageUrl    string `db:"image_url"`
    FriendCount int    `db:"friend_count"`
    CreatedAt   string `db:"created_at"`
}

type PostComment struct {
    Creator   PostUser
}


type PostRawDBResponse struct {
}

type PostFilter struct {
	Limit   int  `json:"limit"`
	Offset  int	 `json:"offset"`
	Search  string  `json:"search"`
	SearchTag  []string `json:"searchTag"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
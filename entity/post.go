package entity

type Post struct {
	ID         string   `json:"postId"`
	UserID     string   `json:"user_id"`
	PostInHtml string   `json:"postInHtml" db:"post_in_html"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt" db:"created_at"`
	UpdatedAt  string   `json:"updatedAt,omitempty"`
}

type PostData struct {
	PostId   string         `json:"postId"`
	Post     PostDetail     `json:"post"`
	Comments []PostComments `json:"comment"`
	Creator  PostUser       `json:"creator"`
}

type PostDetail struct {
	PostInHtml string      `json:"postInHtml"`
	Tags       interface{} `json:"tags"`
	CreatedAt  string      `json:"created_at"`
}

type PostComments struct {
	Comment   string   `json:"comment"`
	Creator   PostUser `json:"creator"`
	CreatedAt string   `json:"createdAt"`
}

type PostUser struct {
	UserID      string `json:"userID"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}

type PostRawDBData struct {
	PostID                    string      `json:"postId" db:"post_id"`
	PostInHTML                string      `json:"postInHtml" db:"post_in_html"`
	Tags                      interface{} `json:"tags" db:"posts_tags"`
	PostCreatedAt             string      `json:"postCreatedAt" db:"posts_created_at"`
	CreatorID                 string      `json:"creatorId" db:"creator_id"`
	CreatorName               string      `json:"creatorName" db:"creator_name"`
	CreatorImageURL           string      `json:"creatorImageUrl" db:"creator_image_url"`
	CreatorFriendCount        int         `json:"creatorFriendCount" db:"creator_friend_count"`
	CreatorCreatedAt          string      `json:"creatorCreatedAt" db:"creator_created_at"`
	Comment                   *string     `json:"comment" db:"comment,omitempty"`
	CommentCreatedAt          *string     `json:"commentCreatedAt" db:"comment_created_at,omitempty"`
	CommentCreatorID          *string     `json:"commentCreatorId" db:"comment_creator_id,omitempty"`
	CommentCreatorName        *string     `json:"commentCreatorName" db:"comment_creator_name,omitempty"`
	CommentCreatorImageUrl    *string     `json:"commentCreatorImageUrl" db:"comment_creator_image_url,omitempty"`
	CommentCreatorFriendCount *int        `json:"commentCreatorFriendCount" db:"comment_creator_friend_count,omitempty"`
}

type PostFilter struct {
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
	Search    string   `json:"search"`
	SearchTag []string `json:"searchTag"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

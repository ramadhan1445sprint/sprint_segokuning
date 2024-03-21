package repo

import (

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type PostRepo interface {
	CreatePost(post *entity.Post) error
	GetPost(userId string) ([]entity.PostData, error)
}

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepo {
	return &postRepo{db}
}

func (r *postRepo) CreatePost(post *entity.Post) error {
	_, err := r.db.Exec("INSERT INTO posts (post_in_html, tags, user_id) VALUES ($1, $2, $3)", post.PostInHtml, post.Tags, post.UserID)

	if err != nil {
		return err
	}

	return nil
}


func (r *postRepo) GetPost(postId string) ([]entity.PostData, error) {
	var posts []entity.PostData

	query := `
		SELECT
			p.id AS "posts.",
			p.post_in_html AS "posts.post_in_html",
			p.tags AS "posts.tags",
			p.created_at AS "posts.created_at",
			u.id AS "creator.id",
			u.name AS "creator.name",
			u.image_url AS "creator.image_url",
			u.friend_count AS "creator.friend_count",
			u.created_at AS "creator.created_at",
			c.comment AS "comments",
			cu.id AS "comments.creator.id",
			cu.name AS "comments.creator.name",
			cu.image_url AS "comments.creator.image_url",
			cu.friend_count AS "comments.creator.friend_count",
			cu.created_at AS "comments.creator.created_at"
		FROM
			posts p
		JOIN
			users u ON p.user_id = u.id
		JOIN
			comments c ON p.id = c.post_id
		JOIN
			users cu ON c.user_id = cu.id
			WHERE p.id = $1
			ORDER BY
		p.created_at DESC, c.created_at DESC
	`

	if err := r.db.QueryRowx(query, postId).Scan(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}

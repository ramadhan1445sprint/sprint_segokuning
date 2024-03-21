package repo

import (

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type CommentRepo interface {
	CreateComment(comment *entity.Comment) error
}

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) CommentRepo {
	return &commentRepo{db}
}

func (r *commentRepo) CreateComment(comment *entity.Comment) error {
	_, err := r.db.Exec("INSERT INTO comments (comment, post_id, user_id) VALUES ($1, $2, $3)", comment.Comment, comment.PostID, comment.UserID)

	if err != nil {
		return err
	}

	return nil
}
package repo

import (

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type CommentRepo interface {
	CreateComment(comment *entity.Comment) error
	CheckPostById(postId string) (bool, error)
	CheckFriendPost(postId string, userId string) (bool, error)
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

func (r *commentRepo) CheckPostById(postId string) (bool, error) {
	var exist int

	err := r.db.Get(&exist, "SELECT count(*) from posts where id = $1", postId)

	if err != nil {
		return false, err
	}

	if exist > 0 {
		return true, nil
	}

	return false, nil

}
func (r *commentRepo) CheckFriendPost(postId string, userId string) (bool, error) {
	var creator string
	var exist int

	err := r.db.Get(&creator, "SELECT user_id from posts where id = $1", postId)
	if err != nil {
		return false, err
	}

	err = r.db.Get(&exist, "SELECT count(*) from friends where (user_id1 = $1 and user_id2 = $2) or (user_id1 = $3 and user_id2 = $4)", userId, creator, creator, userId)
	if err != nil {
		return false, err
	}

	if exist > 0 {
		return true, nil
	}

	return false, nil
}
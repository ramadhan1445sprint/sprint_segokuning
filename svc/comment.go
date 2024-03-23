package svc

import (
	"github.com/google/uuid"
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type CommentSvc interface {
	CreateComment(comment entity.Comment) *customErr.CustomError
}

type commentSvc struct {
	repo repo.CommentRepo
}

func NewCommentSvc(repo repo.CommentRepo) CommentSvc {
	return &commentSvc{repo: repo}
}

func (s *commentSvc) CreateComment(comment entity.Comment) *customErr.CustomError {
	_, err := uuid.Parse(comment.PostID)

	if err != nil {
		custErr := customErr.NewNotFoundError("post not found")
		return &custErr
	}

	exist, err := s.repo.CheckPostById(comment.PostID)

	if err != nil {
		custErr := customErr.NewInternalServerError(err.Error())
		return &custErr
	}

	if !exist {
		custErr := customErr.NewNotFoundError("post not found")
		return &custErr
	}

	if err = s.repo.CreateComment(&comment); err != nil {
		custErr := customErr.NewInternalServerError(err.Error())
		return &custErr
	}

	return nil
}

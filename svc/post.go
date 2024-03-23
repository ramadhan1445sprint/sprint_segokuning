package svc

import (
	"github.com/ramadhan1445sprint/sprint_segokuning/customErr"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
)

type PostSvc interface {
	CreatePost(post entity.Post) *customErr.CustomError
	GetPost(filter entity.PostFilter) (*entity.PostResponse, *customErr.CustomError)
}

type postSvc struct {
	repo repo.PostRepo
}

func NewPostSvc(repo repo.PostRepo) PostSvc {
	return &postSvc{repo}
}

func (s *postSvc) CreatePost(post entity.Post) *customErr.CustomError {

	if err := s.repo.CreatePost(&post); err != nil {
		return &customErr.CustomError{Message: err.Error()}
	}

	return nil
}

func (s *postSvc) GetPost(filter entity.PostFilter) (*entity.PostResponse, *customErr.CustomError) {

	resp, err := s.repo.GetPost(&filter)

	if err != nil {
		custErr := customErr.NewInternalServerError(err.Error())
		return nil, &custErr
	}

	if err != nil {
		custErr := customErr.NewInternalServerError(err.Error())
		return nil, &custErr
	}

	meta := entity.PostMeta{
		Total:  len(resp),
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	postResp := entity.PostResponse{
		Message: "success",
		Data:    resp,
		Meta:    meta,
	}

	return &postResp, nil
}

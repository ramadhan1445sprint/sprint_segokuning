package svc

import "github.com/ramadhan1445sprint/sprint_segokuning/repo"

type Svc interface {}

type svc struct {
	repo repo.Repo
}

func NewSvc(r repo.Repo) Svc {
	return &svc{
		repo: r,
	}
}
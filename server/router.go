package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/controller"
	"github.com/ramadhan1445sprint/sprint_segokuning/middleware"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")

	registerHealthRoute(mainRoute, s.db)
}

func registerHealthRoute(r fiber.Router, db *sqlx.DB) {
	ctr := controller.NewController(svc.NewSvc(repo.NewRepo(db)))

	newRoute(r, "GET", "/health", ctr.HealthCheck)
	newRouteWithAuth(r, "GET", "/auth", ctr.AuthCheck)
}

func newRoute(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.RecordDuration, handler)
}

func newRouteWithAuth(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.RecordDuration, middleware.Auth, handler)
}

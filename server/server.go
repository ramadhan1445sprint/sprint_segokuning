package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ramadhan1445sprint/sprint_segokuning/middleware"
)

type Server struct {
	db  *sqlx.DB
	app *fiber.App
}

func NewServer(db *sqlx.DB) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return &Server{
		db:  db,
		app: app,
	}
}

func (s *Server) Run() error {
	return s.app.Listen(":8080")
}

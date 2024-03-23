package server

import (
	"context"
	"log"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/config"
	"github.com/ramadhan1445sprint/sprint_segokuning/controller"
	"github.com/ramadhan1445sprint/sprint_segokuning/middleware"
	"github.com/ramadhan1445sprint/sprint_segokuning/repo"
	"github.com/ramadhan1445sprint/sprint_segokuning/svc"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")

	registerHealthRoute(mainRoute, s.db)
	registerUserRouter(mainRoute, s.db)
	registerFriendRouter(mainRoute, s.db)
	registerImageRoute(mainRoute)
	registerFriendRouter(mainRoute, s.db)
	registerPostRouter(mainRoute, s.db, s.validator)
}

func registerHealthRoute(r fiber.Router, db *sqlx.DB) {
	ctr := controller.NewController(svc.NewSvc(repo.NewRepo(db)))

	newRoute(r, "GET", "/health", ctr.HealthCheck)
	newRouteWithAuth(r, "GET", "/auth", ctr.AuthCheck)
}

func registerImageRoute(r fiber.Router) {
	bucket := config.GetString("S3_BUCKET_NAME")
	cfg, err := awsCfg.LoadDefaultConfig(
		context.Background(),
		awsCfg.WithRegion("ap-southeast-1"),
		awsCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.GetString("S3_ID"),
				config.GetString("S3_SECRET_KEY"),
				"",
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctr := controller.NewImageController(svc.NewImageSvc(cfg, bucket))

	newRouteWithAuth(r, "POST", "/image", ctr.UploadImage)
}

func registerUserRouter(r fiber.Router, db *sqlx.DB) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))
	userGroup := r.Group("/user")

	newRoute(userGroup, "POST", "/register", ctr.Register)
	newRoute(userGroup, "POST", "/login", ctr.Login)
	newRouteWithAuth(userGroup, "PATCH", "", ctr.UpdateAccountUser)
	newRouteWithAuth(userGroup, "POST", "/link", ctr.UpdateLinkEmailAccount)
	newRouteWithAuth(userGroup, "POST", "/link/phone", ctr.UpdateLinkPhoneAccount)
}

func registerFriendRouter(r fiber.Router, db *sqlx.DB) {
	ctr := controller.NewFriendController(svc.NewFriendSvc(repo.NewUserRepo(db), repo.NewFriendRepo(db)))

	friendGroup := r.Group("/friend")

	newRouteWithAuth(friendGroup, "POST", "/", ctr.AddFriend)
	newRouteWithAuth(friendGroup, "DELETE", "/", ctr.DeleteFriend)
	newRouteWithAuth(friendGroup, "GET", "/", ctr.GetListFriends)
}

func registerPostRouter(r fiber.Router, db *sqlx.DB, validate *validator.Validate) {
	postCtr := controller.NewPostController(svc.NewPostSvc(repo.NewPostRepo(db)), validate)
	commentCtr := controller.NewCommentController(svc.NewCommentSvc(repo.NewCommentRepo(db)), validate)
	postGroup := r.Group("/post")

	newRouteWithAuth(postGroup, "GET", "", postCtr.GetPost)
	newRouteWithAuth(postGroup, "POST", "", postCtr.CreatePost)
	newRouteWithAuth(postGroup, "POST", "/comment", commentCtr.CreateComment)
}

func newRoute(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.RecordDuration, handler)
}

func newRouteWithAuth(router fiber.Router, method, path string, handler fiber.Handler) {
	router.Add(method, path, middleware.RecordDuration, middleware.Auth, handler)
}

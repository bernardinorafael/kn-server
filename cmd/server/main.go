package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormrepo"
	database "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/http/route"
	"github.com/bernardinorafael/kn-server/internal/infra/s3client"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	l := slog.New(logger.NewLog(nil))

	env, err := config.NewConfig()
	if err != nil {
		l.Error("failed load env", "error", err)
		return
	}

	db, err := database.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", "error", err)
		panic(err)
	}

	jwtAuth, err := auth.NewJWTAuth(l, env.JWTSecret)
	if err != nil {
		l.Error("cannot initialize jwt auth")
		panic(err)
	}

	s3, err := s3client.New(env)
	if err != nil {
		l.Error("cannot initialize aws s3 client")
		panic(err)
	}

	userRepo := gormrepo.NewUserRepo(db)
	productRepo := gormrepo.NewProductRepo(db)
	_ = gormrepo.NewTeamRepo(db)

	s3Service := service.NewS3Service(s3, l)
	authService := service.NewAuthService(l, userRepo)
	userService := service.NewUserService(l, userRepo)
	productService := service.NewProductService(service.WithProductParams{
		Log:         l,
		Env:         env,
		ProductRepo: productRepo,
		FileService: s3Service,
	})

	authHandler := route.NewAuthHandler(l, authService, jwtAuth, env)
	productHandler := route.NewProductHandler(l, productService, jwtAuth)
	userHandler := route.NewUserHandler(l, userService, jwtAuth)

	authHandler.RegisterRoute(router)
	productHandler.RegisterRoute(router)
	userHandler.RegisterRoute(router)

	l.Info("server listening", "port", env.Port)

	err = http.ListenAndServe(":"+env.Port, router)
	if err != nil {
		l.Error("error connecting web server", "httperr", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormrepo"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/http/route"
	"github.com/bernardinorafael/kn-server/internal/infra/http/server"
	"github.com/bernardinorafael/kn-server/internal/infra/s3client"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/rs/cors"
)

func main() {
	l := slog.New(logger.NewLog(nil))
	router := server.New()

	env, err := config.NewConfig()
	if err != nil {
		l.Error("failed load env", err.Error(), err)
		return
	}

	con, err := db.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", err.Error(), err)
		panic(err)
	}

	jwtAuth, err := auth.NewJWTAuth(l, env.JWTSecret)
	if err != nil {
		l.Error("cannot initialize jwt auth")
		panic(err)
	}

	s3, err := s3client.New(env)
	if err != nil {
		l.Error("cannot initialize aws s3 clien")
		panic(err)
	}

	userRepo := gormrepo.NewUserRepo(con)
	productRepo := gormrepo.NewProductRepo(con)

	s3Service := service.NewS3Service(s3, env, l)
	authService := service.NewAuthService(l, userRepo)
	productService := service.NewProductService(l, env, productRepo, s3Service)
	userService := service.NewUserService(l, userRepo)

	authHandler := route.NewAuthHandler(l, authService, jwtAuth, env)
	productHandler := route.NewProductHandler(l, productService, jwtAuth)
	userHandler := route.NewUserHandler(userService, jwtAuth)

	authHandler.RegisterRoute(router)
	productHandler.RegisterRoute(router)
	userHandler.RegisterRoute(router)

	l.Info(fmt.Sprintf("server lintening on port %v", env.Port))

	c := cors.Default().Handler(router)
	err = http.ListenAndServe(":"+env.Port, c)
	if err != nil {
		l.Error("error connecting web server", err.Error(), err)
		os.Exit(1)
	}
}

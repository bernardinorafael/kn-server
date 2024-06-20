package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/route"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/server"
	"github.com/rs/cors"
)

func main() {
	l := slog.New(logger.NewLog(nil))
	router := server.New()

	env, err := config.GetConfigEnv()
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

	// Repositories
	userRepo := repository.NewUserRepo(con)
	productRepo := repository.NewProductRepo(con)

	// Services
	authService := service.NewAuthService(l, userRepo)
	productService := service.NewProductService(l, productRepo)
	userService := service.NewUserService(l, userRepo)

	// Handlers
	authHandler := route.NewAuthHandler(l, authService, jwtAuth, env)
	productHandler := route.NewProductHandler(l, productService, jwtAuth)
	userHandler := route.NewUserHandler(userService, jwtAuth)

	// Registering routes
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

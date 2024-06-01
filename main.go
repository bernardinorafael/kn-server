package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/http/handler"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
	"github.com/rs/cors"
)

func main() {
	l := slog.New(logger.NewLog(nil))
	mux := http.NewServeMux()

	env, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed load env", err)
		return
	}

	con, err := db.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", err)
		panic(err)
	}

	/*
	* init repositories
	 */
	userRepo := repository.NewUserRepo(con)
	productRepo := repository.NewProductRepo(con)

	/*
	* init services
	 */
	productService := service.NewProductService(l, productRepo)
	authService := service.NewAuthService(l, userRepo)
	jwtService := service.NewJWTService(l, env)

	/*
	* init handlers
	 */
	authHandler := handler.NewAuthHandler(l, authService, jwtService)
	productHandler := handler.NewProductHandler(l, productService)

	/*
	* init routes
	 */
	authHandler.RegisterRoute(mux)
	productHandler.RegisterRoute(mux)

	l.Info(fmt.Sprintf("server lintening on port %v", env.Port))

	c := cors.Default().Handler(mux)
	err = http.ListenAndServe(":"+env.Port, c)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

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

	userRepo := repository.NewUserRepo(con)
	productRepo := repository.NewProductRepo(con)

	jwtService := service.NewJWTService(l, env)
	authService := service.NewAuthService(l, userRepo)
	productService := service.NewProductService(l, productRepo)

	authHandler := route.NewAuthHandler(l, authService, jwtService)
	productHandler := route.NewProductHandler(l, productService, jwtService)

	authHandler.RegisterRoute(router)
	productHandler.RegisterRoute(router)

	l.Info(fmt.Sprintf("server lintening on port %v", env.Port))

	c := cors.Default().Handler(router)
	err = http.ListenAndServe(":"+env.Port, c)
	if err != nil {
		l.Error("error connecting web server", err.Error(), err)
		os.Exit(1)
	}
}

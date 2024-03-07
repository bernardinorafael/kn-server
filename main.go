package main

import (
	"log/slog"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/database"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/authroute"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/userroute"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	l := slog.New(logger.NewLog(nil))

	cfg, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed to load env", err)
		return
	}

	conn, err := database.Connect(l)
	if err != nil {
		l.Error("error connect database", err)
		return
	}

	// init repositories
	ar := repository.NewAccountRepository(conn)

	// init services
	svc := service.New(
		service.GetLogger(l),
		service.GetConfig(cfg),
		service.GetUserRepository(ar),
	)

	// init handlers
	auth := userroute.NewUserHandler(svc.UserService)
	user := authroute.NewUserHandler(svc.AuthService, svc.JWTService)

	// init routes
	userroute.Start(router, auth, svc.JWTService)
	authroute.Start(router, user)

	_ = router.SetTrustedProxies(nil)
	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		l.Error("error starting server", err)
		os.Exit(1)
	}
}

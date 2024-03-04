package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/database"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/accountroute"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/authroute"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	l := slog.New(logger.NewHandler(&slog.HandlerOptions{
		AddSource: true,
	}))

	cfg, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed to load env", err)
		return
	}

	ctx := context.Background()

	conn, err := database.Connect(ctx, l)
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
		service.GetAccountRepository(ar),
	)

	// init handlers
	auth := accountroute.NewAccountHandler(svc.AccountService)
	user := authroute.NewUserHandler(svc.AuthService, svc.JWTService)

	// init routes
	accountroute.Start(router, auth, svc.JWTService)
	authroute.Start(router, user)

	_ = router.SetTrustedProxies(nil)
	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		l.Error("error starting server", err)
		os.Exit(1)
	}
}

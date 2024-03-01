package main

import (
	"context"
	"log"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/database"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/accountroute"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	cfg, err := config.GetConfigEnv()
	if err != nil {
		log.Fatalf("failed to load env: %s", err)
		return
	}

	ctx := context.Background()
	l := logger.New(cfg)

	if err != nil {
		l.Errorf(ctx, "error initializing auth token: %s", err)
		return
	}

	conn, err := database.Connect(ctx, l)
	if err != nil {
		l.Errorf(ctx, "error connect database: %s", err)
		return
	}

	accountRepository := repository.NewAccountRepository(conn)

	s, err := service.New(
		service.GetAccountRepository(accountRepository),
		service.GetConfig(cfg),
		service.GetLogger(l),
	)
	if err != nil {
		l.Errorf(ctx, "error to get domain services: %s", err)
		return
	}

	accountHandler := accountroute.NewHandler(s.AccountService, s.AuthService)

	accountroute.Start(r, accountHandler)

	_ = r.SetTrustedProxies(nil)
	if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
		l.Fatalf(ctx, "error starting server: %v", err)
	}
}

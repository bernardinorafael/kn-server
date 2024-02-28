package main

import (
	"context"
	"log"

	"github.com/bernardinorafael/gozinho/config"
	"github.com/bernardinorafael/gozinho/config/logger"
	"github.com/bernardinorafael/gozinho/internal/application/service"
	"github.com/bernardinorafael/gozinho/internal/infra/database"
	"github.com/bernardinorafael/gozinho/internal/infra/repository"
	"github.com/bernardinorafael/gozinho/internal/infra/rest/routes/accountroute"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cfg, err := config.GetConfigEnv()
	if err != nil {
		log.Fatalf("failed to load env: %s", err)
		return
	}

	ctx := context.Background()
	l := logger.New(cfg)

	DB, err := database.Connect(ctx, l)
	if err != nil {
		l.Errorf(ctx, "error connect database: %s", err)
		return
	}

	accountRepository := repository.NewAccountRepository(DB)

	services, err := service.New(
		service.GetAccountRepository(accountRepository),
		service.GetLogger(l),
	)
	if err != nil {
		l.Errorf(ctx, "error to get domain services: %s", err)
		return
	}

	accountHandler := accountroute.NewHandler(services.AccountService)
	accountroute.Start(r, accountHandler)

	_ = r.SetTrustedProxies(nil)
	if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
		l.Fatalf(ctx, "error starting server: %v", err)
	}
}

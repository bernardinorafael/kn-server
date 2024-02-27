package main

import (
	"log"
	"log/slog"

	"github.com/bernardinorafael/gozinho/config/env"
	"github.com/bernardinorafael/gozinho/config/logger"
	"github.com/bernardinorafael/gozinho/internal/application/service"
	"github.com/bernardinorafael/gozinho/internal/infra/database"
	"github.com/bernardinorafael/gozinho/internal/infra/repository"
	"github.com/bernardinorafael/gozinho/internal/infra/rest/routes/accountroute"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	logger.InitLogger()

	cfg, err := env.LoadConfig()
	if err != nil {
		slog.Error("failed to load env", err)
		return
	}

	DB, err := database.Connect()
	if err != nil {
		log.Fatalf("error connect database: %v", err)
		return
	}

	ar := repository.NewAccountRepository(DB)

	services, err := service.New(service.GetAccountRepository(ar))
	if err != nil {
		slog.Error("error to get domain services: ", err)
		return
	}

	accountHandler := accountroute.NewHandler(services.AccountService)
	accountroute.Start(router, accountHandler)

	_ = router.SetTrustedProxies(nil)
	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

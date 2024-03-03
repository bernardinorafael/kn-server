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
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/authroute"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	cfg, err := config.GetConfigEnv()
	if err != nil {
		log.Fatalf("failed to load env: %s", err)
		return
	}

	ctx := context.Background()
	l := logger.New(cfg)

	conn, err := database.Connect(ctx, l)
	if err != nil {
		l.Errorf(ctx, "error connect database: %s", err)
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
		l.Fatalf(ctx, "error starting server: %v", err)
	}
}

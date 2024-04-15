package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/database"
	"github.com/bernardinorafael/kn-server/internal/infra/persistence/repository"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/authroute"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/routes/userroute"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	gin.DebugPrintRouteFunc = func(methodFn, pathFn, _ string, _ int) {
		log.Printf("%v %v\n", methodFn, pathFn)
	}

	l := slog.New(logger.NewLog(nil))

	cfg, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed load env", err)
		return
	}

	con, err := database.Connect(l)
	if err != nil {
		l.Error("error connect database", err)
		return
	}

	// init repositories
	userRepository := repository.NewUserRepository(con)

	// init services
	svc := service.New(
		service.GetLogger(l),
		service.GetConfig(cfg),
		service.GetUserRepository(userRepository),
	)

	// init handlers
	auth := userroute.NewUserHandler(svc.UserService)
	user := authroute.NewUserHandler(svc.AuthService, svc.JWTService)

	// init routes
	userroute.Start(router, auth, svc.JWTService)
	authroute.Start(router, user)

	_ = router.SetTrustedProxies(nil)
	err = router.Run("0.0.0.0:" + cfg.Port)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

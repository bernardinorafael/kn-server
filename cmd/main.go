package main

import (
	"log"
	"log/slog"

	"github.com/bernardinorafael/gozinho/config/env"
	"github.com/bernardinorafael/gozinho/config/logger"
	"github.com/bernardinorafael/gozinho/internal/application/controller"
	"github.com/bernardinorafael/gozinho/internal/application/service"
	"github.com/bernardinorafael/gozinho/internal/infra/database"
	"github.com/bernardinorafael/gozinho/internal/infra/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	_, err := env.LoadConfig()
	if err != nil {
		slog.Error("failed to load env", err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		slog.Error("failed to connect database", err)
		return
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	r := gin.Default()

	controller.NewUserController(r, userService)

	if err := r.Run("0.0.0.0:" + env.Env.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

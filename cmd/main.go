package main

import (
	"log"
	"log/slog"

	"github.com/bernardinorafael/gozinho/config/env"
	"github.com/bernardinorafael/gozinho/config/logger"
	"github.com/bernardinorafael/gozinho/internal/application/controllers"
	"github.com/bernardinorafael/gozinho/internal/application/services"
	"github.com/bernardinorafael/gozinho/internal/infra/database"
	"github.com/bernardinorafael/gozinho/internal/infra/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	_, err := env.LoadConfig()
	if err != nil {
		slog.Error("Failed to load env!", err, slog.String("package", "main"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		slog.Error("Failed to connect database!", err)
		return
	}

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	r := gin.Default()

	controllers.NewUserController(r, userService)

	if err := r.Run("0.0.0.0:" + env.Env.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

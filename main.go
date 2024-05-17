package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
)

func main() {
	l := slog.New(logger.NewLog(nil))
	mux := http.NewServeMux()

	env, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed load env", err)
		return
	}

	con, err := db.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", err)
		panic(err)
	}

	userRepo := repository.NewUserRepo(con)

	service.NewAuthService(userRepo, l)

	err = http.ListenAndServe(":"+env.Port, mux)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

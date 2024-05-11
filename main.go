package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
)

func main() {
	router := http.NewServeMux()

	l := slog.New(logger.NewLog(nil))

	env, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed load env", err)
		return
	}

	// init database
	_, err = db.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", err)
		panic(err)
	}

	// init services
	_ = service.New(service.GetLogger(l))

	err = http.ListenAndServe(":"+env.Port, router)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

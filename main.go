package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/config/logger"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	db "github.com/bernardinorafael/kn-server/internal/infra/database/pg"
	routes "github.com/bernardinorafael/kn-server/internal/infra/http"
	"github.com/bernardinorafael/kn-server/internal/infra/repository"
)

func main() {
	l := slog.New(logger.NewLog(nil))
	router := http.NewServeMux()

	router.HandleFunc("GET /user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Write([]byte(id))
	})

	env, err := config.GetConfigEnv()
	if err != nil {
		l.Error("failed load env", err)
		return
	}

	// init database
	con, err := db.Connect(l, env.DSN)
	if err != nil {
		l.Error("error connecting db", err)
		panic(err)
	}

	repository.NewUserRepo(con)

	// init services
	_ = service.New(service.GetLogger(l))

	routes.InitRoutes(router)

	err = http.ListenAndServe(":"+env.Port, router)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

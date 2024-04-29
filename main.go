package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/bernardinorafael/kn-server/config/logger"
)

func main() {
	router := http.NewServeMux()

	// router.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	now := time.Now().Format(time.Kitchen)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)

	// 	w.Write([]byte(now))
	// })

	l := slog.New(logger.NewLog(nil))

	// cfg, err := config.GetConfigEnv()
	// if err != nil {
	// 	l.Error("failed load env", err)
	// 	return
	// }

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		l.Error("error connecting web server", err)
		os.Exit(1)
	}
}

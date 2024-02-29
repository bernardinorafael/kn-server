package logger

import (
	"context"

	"github.com/bernardinorafael/gozinho/config"
	utillog "github.com/bernardinorafael/gozinho/util/log"
)

func New(cfg *config.EnvFile) utillog.Logger {
	params := utillog.LogParams{
		DebugLevel:      cfg.Debug,
		AppName:         cfg.Name,
		AttrFromContext: AttrFromContext,
	}
	return utillog.New(params)
}

func AttrFromContext(ctx context.Context) []any {
	var args []any
	return args
}

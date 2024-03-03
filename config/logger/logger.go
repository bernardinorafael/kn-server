package logger

import (
	"context"

	"github.com/bernardinorafael/kn-server/config"
	utillog "github.com/bernardinorafael/kn-server/helper/log"
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

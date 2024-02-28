package utillog

import (
	"context"
	"log/slog"
)

type LogParams struct {
	slogOptions slog.HandlerOptions

	AppName         string
	DebugLevel      bool
	AttrFromContext func(ctx context.Context) []any
}

func New(params LogParams) Logger {
	return newSlogLogger(params)
}

type Logger interface {
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, msg string, args ...any)
	Infow(ctx context.Context, msg string, keyAndVal ...any)
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, msg string, args ...any)
	Debugw(ctx context.Context, msg string, keyAndVal ...any)
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, msg string, args ...any)
	Warnw(ctx context.Context, msg string, keyAndVal ...any)
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, msg string, args ...any)
	Errorw(ctx context.Context, msg string, keyAndVal ...any)
	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, msg string, args ...any)
	Fatalw(ctx context.Context, msg string, keyAndVal ...any)
	Print(args ...any)
	Printf(msg string, v ...any)
}

package utillog

import (
	"context"
	"log/slog"
)

type LogParams struct {
	slogOptions slog.HandlerOptions

	AppName         string
	DebugLevel      bool
	AttrFromContext func(c context.Context) []any
}

func New(params LogParams) Logger {
	return newSlogLogger(params)
}

type Logger interface {
	Info(c context.Context, msg string)
	Infof(c context.Context, msg string, args ...any)
	Infow(c context.Context, msg string, keyAndVal ...any)
	Debug(c context.Context, msg string)
	Debugf(c context.Context, msg string, args ...any)
	Debugw(c context.Context, msg string, keyAndVal ...any)
	Warn(c context.Context, msg string)
	Warnf(c context.Context, msg string, args ...any)
	Warnw(c context.Context, msg string, keyAndVal ...any)
	Error(c context.Context, msg string)
	Errorf(c context.Context, msg string, args ...any)
	Errorw(c context.Context, msg string, keyAndVal ...any)
	Fatal(c context.Context, msg string)
	Fatalf(c context.Context, msg string, args ...any)
	Fatalw(c context.Context, msg string, keyAndVal ...any)
	Print(args ...any)
	Printf(msg string, v ...any)
}

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

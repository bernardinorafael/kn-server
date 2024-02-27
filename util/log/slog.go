package utillog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

const LevelFatal = "FATAL"
const LevelFatalCode = 60

var CustomLevels = map[int]string{
	LevelFatalCode: LevelFatal,
}

type SlogLogger struct {
	params LogParams
	*slog.Logger
}

func newSlogLogger(params LogParams) *SlogLogger {
	l := &SlogLogger{params: params}

	if params.DebugLevel {
		params.slogOptions.Level = slog.LevelDebug
	}

	l.Logger = slog.New(newLogJSONFormatter(os.Stdout, params))

	return l
}

func (l *SlogLogger) Info(c context.Context, msg string) {
	l.Logger.InfoContext(c, msg, l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Infof(c context.Context, msg string, args ...any) {
	l.Logger.InfoContext(c, fmt.Sprintf(msg, args...), l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Infow(c context.Context, msg string, keyAndVal ...any) {
	l.Logger.InfoContext(c, msg, append(l.params.AttrFromContext(c), keyAndVal...)...)
}

func (l *SlogLogger) Debug(c context.Context, msg string) {
	l.Logger.DebugContext(c, msg, l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Debugf(c context.Context, msg string, args ...any) {
	l.Logger.DebugContext(c, fmt.Sprintf(msg, args...), l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Debugw(c context.Context, msg string, keyAndVal ...any) {
	l.Logger.DebugContext(c, msg, append(l.params.AttrFromContext(c), keyAndVal...)...)
}

func (l *SlogLogger) Warn(c context.Context, msg string) {
	l.Logger.WarnContext(c, msg, l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Warnf(c context.Context, msg string, args ...any) {
	l.Logger.WarnContext(c, fmt.Sprintf(msg, args...), l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Warnw(c context.Context, msg string, keyAndVal ...any) {
	l.Logger.WarnContext(c, msg, append(l.params.AttrFromContext(c), keyAndVal...)...)
}

func (l *SlogLogger) Error(c context.Context, msg string) {
	l.Logger.ErrorContext(c, msg, l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Errorf(c context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(c, fmt.Sprintf(msg, args...), l.params.AttrFromContext(c)...)
}

func (l *SlogLogger) Errorw(c context.Context, msg string, keyAndVal ...any) {
	l.Logger.ErrorContext(c, msg, append(l.params.AttrFromContext(c), keyAndVal...)...)
}

func (l *SlogLogger) Fatal(c context.Context, msg string) {
	l.Logger.Log(c, LevelFatalCode, msg, l.params.AttrFromContext(c)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(c context.Context, msg string, args ...any) {
	l.Logger.Log(c, LevelFatalCode, fmt.Sprintf(msg, args...), l.params.AttrFromContext(c)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalw(c context.Context, msg string, keyAndVal ...any) {
	l.Logger.Log(c, LevelFatalCode, msg, append(l.params.AttrFromContext(c), keyAndVal...)...)
	os.Exit(1)
}

func (l *SlogLogger) Print(args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, "", args...)
}

func (l *SlogLogger) Printf(msg string, args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, fmt.Sprintf(msg, args...))
}

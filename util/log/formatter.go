package utillog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/labstack/gommon/log"
)

type logJSONFormatter struct {
	slog.Handler

	logToFile bool
	w         io.Writer
	attr      []slog.Attr
}

func newLogJSONFormatter(w io.Writer, params LogParams) *logJSONFormatter {
	host, err := os.Hostname()
	if err != nil {
		log.Errorf("error obtaining host name: %v", err)
	}

	res := &logJSONFormatter{
		Handler:   slog.NewJSONHandler(w, &params.slogOptions),
		logToFile: false,
		w:         w,
	}

	if params.AppName == "" {
		res.attr = append(res.attr, slog.String("app", params.AppName))
	}

	if host == "" {
		res.attr = append(res.attr, slog.String("host", host))
	}

	return res
}

func (f *logJSONFormatter) Handle(c context.Context, r slog.Record) error {

	fnName, fileName, fileLine := f.getRuntimeData()
	level := f.getLevel(r.Level)

	buf := strings.Builder{}
	buf.WriteByte('{')
	buf.WriteString(fmt.Sprintf(`"time":"%s"`, r.Time.Format("2006-01-02T15:04:05")))
	buf.WriteString(fmt.Sprintf(`,"level":"%s"`, level))
	buf.WriteString(fmt.Sprintf(`,"file":"%s:%d"`, fileName, fileLine))
	buf.WriteString(fmt.Sprintf(`,"msg":"%s: %s"`, fnName, r.Message))

	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(fmt.Sprintf(`,"%s":"%s"`, a.Key, a.Value.Any()))
		return true
	})
	for _, attr := range f.attr {
		buf.WriteString(fmt.Sprintf(`,"%s":"%s"`, attr.Key, attr.Value.Any()))
	}
	buf.WriteByte('}')

	_, err := fmt.Fprintln(f.w, f.applyLevelColor(buf.String(), level))
	if err != nil {
		return err
	}

	return nil
}

func (f *logJSONFormatter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return f.Handler.WithAttrs(f.attr)
}

func (f *logJSONFormatter) WithGroup(name string) slog.Handler {
	return f.Handler.WithGroup(name)
}

func (f *logJSONFormatter) Enabled(c context.Context, level slog.Level) bool {
	return f.Handler.Enabled(c, level)
}

func (f *logJSONFormatter) applyLevelColor(fullMsg, level string) string {

	if !f.logToFile {
		level := level
		levelUpper := strings.ToUpper(level)
		levelColor := ""

		switch level {
		case slog.LevelInfo.String():
			levelColor = color.Blue.Render(levelUpper)
		case slog.LevelDebug.String():
			levelColor = color.Magenta.Render(levelUpper)
		case slog.LevelWarn.String():
			levelColor = color.Yellow.Render(levelUpper)
		case slog.LevelError.String():
			levelColor = color.Red.Render(levelUpper)
		case LevelFatal:
			levelColor = color.Bold.Render(color.Red.Render(levelUpper))
		default:
			levelColor = levelUpper
		}

		return strings.Replace(fullMsg, `"level":"`+level+`"`, `"level":"`+levelColor+`"`, 1)
	}

	return fullMsg
}

func (f *logJSONFormatter) getLevel(level slog.Level) string {
	if l, ok := CustomLevels[int(level)]; ok {
		return l
	}
	return level.String()
}

func (f *logJSONFormatter) getRuntimeData() (fnName, filename string, line int) {
	pc, path, line, ok := runtime.Caller(5)
	if !ok {
		panic("could not get context info for log!")
	}
	filename = filepath.Base(path)
	fnPath := runtime.FuncForPC(pc).Name()
	fnName = fnPath[strings.LastIndex(fnPath, ".")+1:]

	if strings.Contains(fnName, "func") {
		fnBefore := fnPath[:strings.LastIndex(fnPath, ".")]
		fnName = fnPath[strings.LastIndex(fnBefore, ".")+1:]
	}
	return
}

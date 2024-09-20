package logger

import (
	"context"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

/* ------------------------------------------------  自定义handler  ------------------------------------------------ */

type CustomHandler struct {
	handler slog.Handler
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	if requestID, ok := ctx.Value("X-Request-ID").(string); ok {
		record.AddAttrs(slog.String("X-Request-ID", requestID))
	}
	return h.handler.Handle(ctx, record)
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

/* ------------------------------------------------  初始化  ------------------------------------------------ */

var baseLogger *slog.Logger

func Init(file string, level string, tags map[any]any) {
	var output io.Writer
	output = os.Stdout
	if file != "" {
		output = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    20, // 以MB为单位
			MaxBackups: 5,
			MaxAge:     30,   // 以天为单位
			Compress:   true, // 压缩旧日志文件
		}
	}

	levelMap := map[string]slog.Leveler{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	jsonHandler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		AddSource: false,
		Level:     levelMap[level],
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				_, offset := t.Zone()
				return slog.String(a.Key, t.Format(fmt.Sprintf("2006-01-02 15:04:05.000 UTC+%d", offset/3600)))
			}
			return a
		},
	})

	baseLogger = slog.New(&CustomHandler{jsonHandler})
	for k, v := range tags {
		baseLogger = baseLogger.With(k, v)
	}
}

/* ------------------------------------------------  调用堆栈 ------------------------------------------------ */

func caller() string {
	_, file, line, _ := runtime.Caller(2)
	pathSlice := strings.Split(file, "/")
	if len(pathSlice) < 3 {
		return fmt.Sprintf("%s:%d", file, line)
	}
	path := strings.Join(pathSlice[len(pathSlice)-2:], "/")
	return fmt.Sprintf("%s:%d", path, line)
}

/* ------------------------------------------------  baseLogger ------------------------------------------------ */

func DebugContext(ctx context.Context, msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.ErrorContext(ctx, msg, args...)
}

func Debug(msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	args = append(args, "caller", caller())
	baseLogger.Error(msg, args...)
}

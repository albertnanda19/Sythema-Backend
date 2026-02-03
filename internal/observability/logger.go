package observability

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"synthema/internal/config"
)

type Logger struct {
	std *slog.Logger
}

func NewLogger(cfg config.Config) *Logger {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLevel(cfg.LogLevel)})
	l := slog.New(h).With(
		slog.String("app", cfg.AppName),
		slog.String("env", cfg.Environment),
	)
	return &Logger{std: l}
}

func (l *Logger) Info(msg string) {
	l.std.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.std.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.std.Error(msg)
}

func (l *Logger) InfoContext(ctx context.Context, msg string) {
	l.std.InfoContext(ctx, msg)
}

func (l *Logger) WarnContext(ctx context.Context, msg string) {
	l.std.WarnContext(ctx, msg)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string) {
	l.std.ErrorContext(ctx, msg)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "error":
		return slog.LevelError
	case "warn", "warning":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

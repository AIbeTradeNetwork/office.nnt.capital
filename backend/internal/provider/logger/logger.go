package logger

import (
	"sync"

	"server/internal/config"
)

type Logger struct {
	Level string
}

var (
	service Logger
	once    sync.Once
)

// Get logger
func Get() *Logger {
	once.Do(func() {
		cfg := config.Get()
		service = Logger{Level: cfg.LogLevel}
	})
	return &service
}

func (s *Logger) Info(args ...any) {
	//slog.Info("SLOG", args)
}

func (s *Logger) Warn(args ...any) {
	//slog.Warn("SLOG", args)
}

func (s *Logger) Error(args ...any) {
	//slog.Error("SLOG", args)
}

func (s *Logger) Debug(args ...any) {
	//slog.Debug("SLOG", args)
}

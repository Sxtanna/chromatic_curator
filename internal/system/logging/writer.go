package logging

import (
	"io"
	"log/slog"
	"strings"
)

type slogWriter struct {
	logger *slog.Logger
}

func NewSlogWriter(logger *slog.Logger) io.Writer {
	return &slogWriter{
		logger: logger,
	}
}

func (w *slogWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSuffix(string(p), "\n")

	switch {
	case strings.HasPrefix(msg, "ERROR:"):
		w.logger.Error(strings.TrimPrefix(msg, "ERROR:"))
	case strings.HasPrefix(msg, "WARN:"):
		w.logger.Warn(strings.TrimPrefix(msg, "WARN:"))
	case strings.HasPrefix(msg, "INFO:"):
		w.logger.Info(strings.TrimPrefix(msg, "INFO:"))
	case strings.HasPrefix(msg, "DEBUG:"):
		w.logger.Debug(strings.TrimPrefix(msg, "DEBUG:"))
	default:
		w.logger.Info(msg)
	}

	return len(p), nil
}

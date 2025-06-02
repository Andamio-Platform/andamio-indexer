package logutils

import (
	"log/slog"
)

// SlogWriter is a custom io.Writer that writes to a slog.Logger.
type SlogWriter struct {
	logger *slog.Logger
}

// NewSlogWriter creates a new SlogWriter.
func NewSlogWriter(logger *slog.Logger) *SlogWriter {
	return &SlogWriter{logger: logger}
}

// Write implements the io.Writer interface for SlogWriter.
func (w *SlogWriter) Write(p []byte) (n int, err error) {
	w.logger.Info("fiber_log", "message", string(p))
	return len(p), nil
}

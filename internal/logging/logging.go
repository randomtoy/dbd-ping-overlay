// Package logging configures the application's structured logger.
package logging

import (
	"io"
	"log/slog"
)

// New returns a slog.Logger that writes human-readable text records to w.
// It is intended for development use, where logs are read directly from a
// console window.
func New(w io.Writer) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}

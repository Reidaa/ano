package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func NewJSON() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

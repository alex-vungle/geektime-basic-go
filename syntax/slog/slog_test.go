package slog

import (
	"log/slog"
	"runtime"
	"testing"
)

func TestSlog(t *testing.T) {
	slog.Info("打印 debug", slog.Int64("id", 123))
	l := slog.With(slog.Int64("id", 123))
	l.Info("打印 with")
	l.Error("error with")

	runtime.NumGoroutine()
}

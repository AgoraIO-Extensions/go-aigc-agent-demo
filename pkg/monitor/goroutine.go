package monitor

import (
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"runtime"
)

func LogGoroutines() {
	logger.Info("[engine] final goroutines nums", slog.Int("num", runtime.NumGoroutine()))
}

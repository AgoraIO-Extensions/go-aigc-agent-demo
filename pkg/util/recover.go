package util

import (
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
	"runtime"
)

func Recover() {
	if r := recover(); r != nil {
		logger.Error(fmt.Sprintf("Recovered from panic: %v", r))
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		logger.Error(fmt.Sprintf("Stack trace:\n%s", buf[:stackSize]))
	}
}

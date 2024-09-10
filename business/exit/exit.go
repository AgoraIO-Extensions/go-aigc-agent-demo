package exit

import (
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"os"
	"time"
)

type ExitManager struct {
	StartTime   int64
	MaxLifeTime int64 // Maximum process uptime.
}

func NewExitManager(startTime int64, maxLifeTime int64) *ExitManager {
	return &ExitManager{
		StartTime:   startTime,
		MaxLifeTime: maxLifeTime,
	}
}

// OnUserLeft Handle user departure events (currently supports only single user scenarios)
func (e *ExitManager) OnUserLeft(conn *agoraservice.RtcConnection, uid string, reason int) {
	logger.Info("[exit] User has left; the process is about to exit.", slog.String("uid", uid))
	os.Exit(0)
}

func (e *ExitManager) HandlerMaxLifeTime() {
	go func() {
		leftLifeTime := e.MaxLifeTime - (time.Now().Unix() - e.StartTime)
		if leftLifeTime <= 0 {
			logger.Info("Reached maximum uptime; exiting soon...")
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("Remaining uptime: %d", leftLifeTime))
		<-time.After(time.Second * time.Duration(leftLifeTime))
		logger.Info("Reached maximum uptime; exiting soon...")
		os.Exit(0)
	}()
}

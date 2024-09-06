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
	MaxLifeTime int64 // 进程最大存活时间
}

func NewExitManager(startTime int64, maxLifeTime int64) *ExitManager {
	return &ExitManager{
		StartTime:   startTime,
		MaxLifeTime: maxLifeTime,
	}
}

// OnUserLeft 处理用户离开事件（当前只支持单个用户的情况）
func (e *ExitManager) OnUserLeft(conn *agoraservice.RtcConnection, uid string, reason int) {
	logger.Info("[exit] 用户离开，进程即将退出", slog.String("uid", uid))
	os.Exit(0)
}

func (e *ExitManager) HandlerMaxLifeTime() {
	go func() {
		leftLifeTime := e.MaxLifeTime - (time.Now().Unix() - e.StartTime)
		if leftLifeTime <= 0 {
			logger.Info("达到最大存活时间，即将退出...")
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("剩余存活时间：%d", leftLifeTime))
		<-time.After(time.Second * time.Duration(leftLifeTime))
		logger.Info("达到最大存活时间，即将退出...")
		os.Exit(0)
	}()
}

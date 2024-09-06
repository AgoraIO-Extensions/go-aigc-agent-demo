package ali

import (
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type STT struct {
	SID  int64
	conn *conn
}

func NewSTT(sid int64, cfg *Config) (*STT, error) {
	ist := &STT{
		SID: sid,
	}
	ist.conn = cfg.connPool.GetConn()
	ist.conn.sid = sid
	return ist, nil
}

func (stt *STT) Send(chunk []byte, end bool) error {
	if end {
		ready, err := stt.conn.nlsST.Stop() // 通知服务端语音发送完毕，之后服务端会回调 onSentenceEnd 函数
		logger.Info("[stt] 停止推流", slog.Int64("sid", stt.SID))
		if err != nil {
			stt.conn.nlsST.Shutdown() // 强制停止实时语音识别
			return fmt.Errorf("[nlsST.Stop]%v", err)
		}
		go func(sid int64) {
			if er := waitReady(ready); er != nil {
				logger.Error(fmt.Sprintf("[Stop waitReady]%v", er), slog.Int64("sid", sid))
			}
			stt.conn.nlsST.Shutdown()
		}(stt.SID)
		return nil
	}

	if err := stt.conn.nlsST.SendAudioData(chunk); err != nil {
		return err
	}
	return nil
}

func (stt *STT) GetResult() <-chan *common.Result {
	return stt.conn.result
}

package ali

import (
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
)

type InternalSTT struct {
	SID  int64
	conn *conn
}

func NewInternalSTT(sid int64, cfg *Config) (*InternalSTT, error) {
	ist := &InternalSTT{
		SID: sid,
	}
	ist.conn = cfg.connPool.GetConn()
	ist.conn.sid = sid
	return ist, nil
}

func (ist *InternalSTT) Send(chunk []byte, end bool) error {
	if end {
		ready, err := ist.conn.nlsST.Stop() // 通知服务端语音发送完毕，之后服务端会回调 onSentenceEnd 函数
		logger.Inst().Info("[stt] 停止推流", zap.Int64("sid", ist.SID))
		if err != nil {
			ist.conn.nlsST.Shutdown() // 强制停止实时语音识别
			return fmt.Errorf("[nlsST.Stop]%v", err)
		}
		go func(sid int64) {
			if er := waitReady(ready); er != nil {
				logger.Inst().Error(fmt.Sprintf("[Stop waitReady]%v", er), zap.Int64("sid", sid))
			}
			ist.conn.nlsST.Shutdown()
		}(ist.SID)
		return nil
	}

	if err := ist.conn.nlsST.SendAudioData(chunk); err != nil {
		return err
	}
	return nil
}

func (ist *InternalSTT) GetResult() <-chan *common.Result {
	return ist.conn.result
}

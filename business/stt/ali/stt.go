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
		ready, err := stt.conn.nlsST.Stop() // Notify the server that the voice has been sent; thereafter, the server will call the onSentenceEnd function.
		logger.Info("[stt] Stop pushing stream", slog.Int64("sid", stt.SID))
		if err != nil {
			stt.conn.nlsST.Shutdown() // Forcefully stop real-time speech recognition.
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

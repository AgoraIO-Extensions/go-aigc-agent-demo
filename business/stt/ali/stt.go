package ali

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
)

type STT struct {
	ctx  context.Context
	conn *conn
}

func NewSTT(ctx context.Context, cfg *Config) (*STT, error) {
	ist := &STT{
		ctx: ctx,
	}
	ist.conn = cfg.connPool.GetConn()
	ist.conn.ctx = ctx // ist.conn is created asynchronously, and ist.conn.ctx is ‘context.Background()’, so the ctx needs to be set here
	return ist, nil
}

func (stt *STT) Send(chunk []byte, end bool) error {
	if end {
		ready, err := stt.conn.nlsST.Stop() // Notify the server that the voice has been sent; thereafter, the server will call the onSentenceEnd function.
		logger.InfoContext(stt.ctx, "[stt] Stop pushing stream")
		if err != nil {
			stt.conn.nlsST.Shutdown() // Forcefully stop real-time speech recognition.
			return fmt.Errorf("[nlsST.Stop]%v", err)
		}
		go func() {
			if er := waitReady(ready); er != nil {
				logger.ErrorContext(stt.ctx, fmt.Sprintf("[Stop waitReady]%v", er))
			}
			stt.conn.nlsST.Shutdown()
		}()
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

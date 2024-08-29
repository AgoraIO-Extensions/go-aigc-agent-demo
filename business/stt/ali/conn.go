package ali

import (
	"fmt"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"go-aigc-agent-demo/business/stt/common"
	"time"
)

type conn struct {
	nlsST   *nls.SpeechTranscription
	expTime time.Time // 连接过期截止时间 10s（取8s）

	sid         int64
	result      chan *common.Result
	returnedAns bool
}

func newConn(cfg *Config, sid int64) (*conn, error) {
	c := &conn{
		expTime: time.Now().Add(time.Second * 8),
		sid:     sid,
		result:  make(chan *common.Result, 100),
	}

	nlsST, err := nls.NewSpeechTranscription(cfg.nlsConnConf, cfg.nlsLogger, c.onTaskFailed, c.onStarted, c.onSentenceBegin, c.onSentenceEnd, c.onResultChanged, c.onCompleted, c.onClose, cfg.nlsLogger)
	if err != nil {
		return nil, fmt.Errorf("[nls.NewSpeechTranscription]%v", err)
	}
	ready, err := nlsST.Start(cfg.nlsStartParamsConf, nil) // 建立连接
	if err != nil {
		nlsST.Shutdown()
		return nil, fmt.Errorf("[nlsST.Start]%v", err)
	}
	err = waitReady(ready)
	if err != nil {
		nlsST.Shutdown()
		return nil, fmt.Errorf("[waitReady]%v", err)
	}
	c.nlsST = nlsST

	return c, nil
}

func waitReady(ch chan bool) error {
	select {
	case done := <-ch:
		{
			if !done {
				return fmt.Errorf("wait failed")
			}
		}
	case <-time.After(5 * time.Second):
		{
			return fmt.Errorf("wait timeout")
		}
	}
	return nil
}

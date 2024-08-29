package stt

import (
	"fmt"
	"go-aigc-agent-demo/business/stt/ali"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/alibaba/speech"
)

type InternalSTT interface {
	Send(chunk []byte, end bool) error
	GetResult() <-chan *common.Result
}

type STT struct {
	vendorName config.SttSelect
	aliConfig  *ali.Config
}

func NewSTT(vendorName config.SttSelect, sttConfig config.STT) (*STT, error) {
	stt := &STT{vendorName: vendorName}
	var err error
	switch vendorName {
	case config.AliSTT:
		c := sttConfig.Ali
		stt.aliConfig, err = ali.Init(c.URL, c.AppKey, speech.TOKEN)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("vendorname 传值不符合预期")
	}
	return stt, nil
}

func (s *STT) GetOneConnection(sid int64) (InternalSTT, error) {
	switch s.vendorName {
	case config.AliSTT:
		ist, err := ali.NewInternalSTT(sid, s.aliConfig)
		if err != nil {
			return nil, fmt.Errorf("[ali.NewInternalSTT]%v", err)
		}
		return ist, nil
	default:
		return nil, fmt.Errorf("vendorname错误")
	}
}

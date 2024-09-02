package tts

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/tts/vendors/ali"
	"go-aigc-agent-demo/config"
)

type TTS struct {
	httpSender *httpSender
}

func Init() (*TTS, error) {
	cfg := config.Inst()
	t := &TTS{}

	switch cfg.TTS.Select {
	case config.AliTTS:
		t.httpSender = &httpSender{client: ali.NewAliTTS(), concurrence: make(chan struct{}, 2)}
	default:
		return nil, fmt.Errorf("tts vendor选择错误")
	}

	return t, nil
}

func (t *TTS) NewSender(sid int64) *Sender {
	s := &Sender{
		sid:        sid,
		httpSender: t.httpSender,
		sentence: &Sentence{
			ID:        sid,
			AudioChan: make(chan []byte, 1000),
		},
	}
	if t.httpSender != nil {
		s.sentence.segChan = make(chan *Segment, 1000)
		s.sentence.mergeSegments() // 异步获取并合并音频结果
	}
	return s
}

type Sender struct {
	sid        int64
	httpSender *httpSender
	sentence   *Sentence
}

// Send text为 "" 则表示当前session发送结束
func (s *Sender) Send(ctx context.Context, segID int, text string) {
	s.httpSender.send(ctx, s.sentence, segID, text)
}

func (s *Sender) GetResult() <-chan []byte {
	return s.sentence.AudioChan
}

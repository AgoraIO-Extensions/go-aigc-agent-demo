package tts

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/tts/vendors/ali"
	"go-aigc-agent-demo/config"
)

type TTS struct {
	httpSender      *httpSender
	websocketSender *webSocketSender
}

func Init() (*TTS, error) {
	cfg := config.Inst()
	t := &TTS{}

	switch cfg.TTS.Select {
	case config.AliTTS:
		t.httpSender = &httpSender{client: ali.NewAliTTS(), concurrence: make(chan struct{}, 2)}
	}

	return t, nil
}

type Sender struct {
	sid             int64
	httpSender      *httpSender
	websocketSender *webSocketSender
	sentence        *Sentence
}

func (t *TTS) NewSender(sid int64) *Sender {
	s := &Sender{
		sid:             sid,
		httpSender:      t.httpSender,
		websocketSender: t.websocketSender,
		sentence: &Sentence{
			ID:        sid,
			AudioChan: make(chan []byte, 1000),
		},
	}
	if t.httpSender != nil {
		s.sentence.segChan = make(chan *Segment, 1000)
		s.sentence.mergeSegments()
	}
	return s
}

// Send text为 "" 则表示当前session发送结束
func (s *Sender) Send(ctx context.Context, segID int, text string) error {
	if s.httpSender != nil {
		s.httpSender.send(ctx, s.sentence, segID, text)
		return nil
	}
	if err := s.websocketSender.send(s.sentence, text); err != nil {
		return fmt.Errorf("websocketSender.send]%v", err)
	}
	return nil
}

func (s *Sender) GetResult() <-chan []byte {
	return s.sentence.AudioChan
}

package common

import (
	"context"
	"errors"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/pkg/logger"
	"io"
	"log/slog"
	"time"
)

type StreamAsk func(ctx context.Context, text string) (io.ReadCloser, error)

type HttpSender struct {
	askFunc     StreamAsk
	concurrence chan struct{}
	sentence    *Sentence
}

func NewHttpSender(sid int64, con int, askFunc StreamAsk) *HttpSender {
	sender := &HttpSender{
		askFunc:     askFunc,
		concurrence: make(chan struct{}, con),
		sentence: &Sentence{
			ID:        sid,
			SegChan:   make(chan *Segment, 1000),
			AudioChan: make(chan []byte, 1000),
		},
	}
	sender.sentence.mergeSegments() // 异步地合并 Sentence下各个segment的audio
	return sender
}

// Send 将segment同步入队后，再异步并发地请求tts
func (h *HttpSender) Send(ctx context.Context, segID int, text string) {
	if text == "" {
		close(h.sentence.SegChan)
		return
	}

	seg := &Segment{
		AudioChan: make(chan []byte, 1000), // 10ms * 1000 = 10s (320KB)
		Sid:       h.sentence.ID,
		ID:        segID,
		Text:      text,
	}
	h.sentence.SegChan <- seg
	h.concurrence <- struct{}{}
	go func(seg *Segment) {
		defer func() {
			<-h.concurrence
		}()
		h.sendSeg(ctx, seg)
	}(seg)
	return
}

func (h *HttpSender) sendSeg(ctx context.Context, seg *Segment) {
	defer close(seg.AudioChan)

	seg.SendTime = time.Now()
	logger.Debug("[tts]发送segment给tts", slog.Int64("sid", seg.Sid), slog.String("seg", seg.Text), slog.Int("seg_id", seg.ID))
	rc, err := h.askFunc(ctx, seg.Text)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("[tts] http请求被打断", slog.String("msg", err.Error()))
			return
		}
		logger.Error("[tts] http流式请求失败", slog.Any("err", err))
		return
	}
	defer rc.Close()
	logger.Debug("[tts]收到响应头和status_code", slog.Int64("sid", seg.Sid), slog.String("seg", seg.Text), slog.Int("seg_id", seg.ID))

	var (
		buf         = make([]byte, 320)
		alreadyRead = 0
		chunkIndex  = 0
	)
	for {
		n, err := rc.Read(buf[alreadyRead:])
		if err == io.EOF {
			return
		}
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("[tts] 读取http返回流被打断", slog.String("msg", err.Error()))
				return
			}
			logger.Error("[tts] 读取http返回流失败", slog.Any("err", err))
			return
		}
		alreadyRead += n
		if alreadyRead < 320 {
			continue
		}
		if chunkIndex == 0 {
			segDur := time.Since(seg.SendTime)
			logger.Info("[tts]<duration>收到一个segment中的首个chunk", sentencelifecycle.Tag(seg.Sid), slog.Int("seg_id", seg.ID),
				slog.String("seg", seg.Text), slog.Int64("dur", segDur.Milliseconds()))
		}
		seg.AudioChan <- buf
		chunkIndex++
		alreadyRead = 0
		buf = make([]byte, 320)
	}
}

func (h *HttpSender) Result() <-chan []byte {
	return h.sentence.AudioChan
}

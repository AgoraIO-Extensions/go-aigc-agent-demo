package tts

import (
	"context"
	"errors"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"io"
	"time"
)

type HttpTTS interface {
	StreamAsk(ctx context.Context, text string) (io.ReadCloser, error)
}

type httpSender struct {
	client      HttpTTS
	concurrence chan struct{}
}

// send 将segment同步入队后，再异步并发地请求tts
func (h *httpSender) send(ctx context.Context, ss *Sentence, segID int, text string) {
	if text == "" {
		close(ss.segChan)
		return
	}

	seg := &Segment{
		AudioChan: make(chan []byte, 1000), // 10ms * 1000 = 10s (320KB)
		Sid:       ss.ID,
		ID:        segID,
		Text:      text,
	}
	ss.segChan <- seg
	h.concurrence <- struct{}{}
	go func(seg *Segment) {
		defer func() {
			<-h.concurrence
		}()
		h.sendSeg(ctx, seg)
	}(seg)
	return
}

func (h *httpSender) sendSeg(ctx context.Context, seg *Segment) {
	defer close(seg.AudioChan)

	seg.SendTime = time.Now()
	logger.Inst().Debug("[tts]发送segment给tts", zap.Int64("sid", seg.Sid), zap.String("seg", seg.Text), zap.Int("seg_id", seg.ID))
	rc, err := h.client.StreamAsk(ctx, seg.Text)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Inst().Info("[tts] http请求被打断", zap.String("msg", err.Error()))
			return
		}
		logger.Inst().Error("[tts] http流式请求失败", zap.Error(err))
		return
	}
	defer rc.Close()
	logger.Inst().Debug("[tts]收到响应头和status_code", zap.Int64("sid", seg.Sid), zap.String("seg", seg.Text), zap.Int("seg_id", seg.ID))

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
				logger.Inst().Info("[tts] 读取http返回流被打断", zap.String("msg", err.Error()))
				return
			}
			logger.Inst().Error("[tts] 读取http返回流失败", zap.Error(err))
			return
		}
		alreadyRead += n
		if alreadyRead < 320 {
			continue
		}
		if chunkIndex == 0 {
			seg.RevTime = time.Now()
			segDur := seg.RevTime.Sub(seg.SendTime)
			logger.Inst().Info("[tts]<duration>收到一个segment中的首个chunk", sentencelifecycle.Tag(seg.Sid), zap.Int("seg_id", seg.ID),
				zap.String("seg", seg.Text), zap.Int64("dur", segDur.Milliseconds()))
		}
		seg.AudioChan <- buf
		chunkIndex++
		alreadyRead = 0
		buf = make([]byte, 320)
	}
}

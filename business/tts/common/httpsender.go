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
	sender.sentence.mergeSegments() // Asynchronously merge the audio of each segment under the Sentence
	return sender
}

// Send After synchronously enqueuing the segment, asynchronously and concurrently request TTS
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
	logger.Info("[tts] Send segment to TTS", slog.Int64("sid", seg.Sid), slog.String("seg", seg.Text), slog.Int("seg_id", seg.ID))
	rc, err := h.askFunc(ctx, seg.Text)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("[tts] HTTP request was interrupted", slog.String("msg", err.Error()))
			return
		}
		logger.Error("[tts] Failed to request tts", slog.Any("err", err))
		return
	}
	defer rc.Close()
	logger.Info("[tts]Get HTTP response header and statusCode", slog.Int64("sid", seg.Sid), slog.String("seg", seg.Text), slog.Int("seg_id", seg.ID))

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
				logger.Info("[tts] Reading the HTTP response stream was interrupted", slog.String("msg", err.Error()))
				return
			}
			logger.Error("[tts] Failed to read HTTP response stream", slog.Any("err", err))
			return
		}
		alreadyRead += n
		if alreadyRead < 320 {
			continue
		}
		if chunkIndex == 0 {
			segDur := time.Since(seg.SendTime)
			logger.Info("[tts]<duration> Received the first chunk of a segment", sentencelifecycle.Tag(seg.Sid), slog.Int("seg_id", seg.ID),
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
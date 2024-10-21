package engine

import (
	"go-aigc-agent-demo/business/aigcCtx"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

func (e *Engine) ProcessSendRTC(input <-chan *ttsResult) {
	for {
		r := <-input
		ctx := r.ctx
		e.sendAudioToRTC(ctx, r.audio)
	}
}

func (e *Engine) sendAudioToRTC(ctx *aigcCtx.AIGCContext, audioChan <-chan []byte) {
	firstSend := true
	for {
		var chunk []byte
		var ok bool
		select {
		case chunk, ok = <-audioChan:
			if !ok {
				logger.InfoContext(ctx, "[rtc] Completed sending audio to RTC.")
				return
			}
			break
		case <-ctx.Done():
			logger.InfoContext(ctx, "[rtc] Interrupted while sending audio to RTC.")
			return
		}
		if firstSend {
			firstSend = false
			ctx.MetaData.StageSendToRTC = true
			logger.InfoContext(ctx, "[rtc] Started sending audio to RTC.")
			logger.InfoContext(ctx, "[sentence]<duration> filter output the tail chunk ——> send the head chunk to RTC", slog.Int64("dur", time.Since(ctx.MetaData.FilterAudioTailRcvTime).Milliseconds()))
		}
		if err := e.rtc.SendPcm(chunk); err != nil {
			logger.ErrorContext(ctx, "[rtc] Failed to send audio to RTC.", slog.Any("err", err))
		}
	}
}

package engine

import (
	"context"
	"go-aigc-agent-demo/business/sentence"
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

func (e *Engine) sendAudioToRTC(ctx context.Context, audioChan <-chan []byte) {
	firstSend := true
quickSend:
	for i := 0; i < 18; i++ { // The instantaneous limit for sending packets is 18 packets; otherwise, the average packet rate must be maintained at 1 packet per 10ms.
		var chunk []byte
		var ok bool
		select {
		case chunk, ok = <-audioChan:
			break
		case <-ctx.Done():
			logger.InfoContext(ctx, "[rtc] Interrupted while sending audio to RTC.")
			return
		}
		if !ok {
			logger.InfoContext(ctx, "[rtc] Completed sending audio to RTC.")
			return
		}
		if firstSend {
			firstSend = false
			sMetaData := sentence.GetMetaData(ctx)
			sMetaData.StageSendToRTC = true
			logger.InfoContext(ctx, "[rtc] Started sending audio to RTC.")
			logger.InfoContext(ctx, "[sentence]<duration> filter output the tail chunk ——> send the head chunk to RTC", slog.Int64("dur", time.Since(sMetaData.FilterAudioTailRcvTime).Milliseconds()))
		}

		if err := e.rtc.SendPcm(chunk); err != nil {
			logger.ErrorContext(ctx, "[rtc] Failed to send audio to RTC.", slog.Any("err", err))
			return
		}
	}

	firstSendTime := time.Now()
	sendCount := 0
	shouldSendCount := 0
	for {
		time.Sleep(time.Millisecond * 50)
		shouldSendCount = int(time.Since(firstSendTime).Milliseconds())/10 - sendCount
		if shouldSendCount > 18 { // If the operation below (<-audioChan) is blocked for a long time (>=140ms), then shouldSendCount will be greater than 18：(140+50)/10=19>18
			logger.InfoContext(ctx, "[rtc] The blocking time is too long; executing quickSend.")
			goto quickSend
		}
		for i := 0; i < shouldSendCount; i++ {
			var readStart = time.Now()
			var chunk []byte
			var ok bool
			select {
			case chunk, ok = <-audioChan:
				break
			case <-ctx.Done():
				logger.InfoContext(ctx, "[rtc] Interrupted while sending audio to RTC.")
				return
			}
			if !ok {
				logger.InfoContext(ctx, "[rtc] Audio sent to RTC completed")
				return
			}

			if dur := time.Since(readStart); dur > time.Millisecond*10 {
				logger.WarnContext(ctx, "[rtc] While sending audio to RTC, blocking in audio retrieval for more than 10ms.", slog.Int64("dur", dur.Milliseconds()))
			}
			if err := e.rtc.SendPcm(chunk); err != nil {
				logger.ErrorContext(ctx, "[rtc] Failed to send audio to RTC.", slog.Any("err", err))
				return
			}
			sendCount++
		}
		if shouldSendCount == 18 {
			firstSendTime = time.Now()
			sendCount = 0
		}
	}
}

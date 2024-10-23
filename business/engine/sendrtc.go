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
			if waitSubsequentNodes(ctx) {
				return
			}
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

func waitSubsequentNodes(ctx *aigcCtx.AIGCContext) bool {
	if dur := time.Since(ctx.MetaData.FilterAudioTailRcvTime); dur.Milliseconds() < 800 {
		waitTime := time.Millisecond*800 - dur
		logger.InfoContext(ctx, "[rtc]<duration> waiting for the subsequent nodes come.", slog.Int64("dur", waitTime.Milliseconds()))
		time.Sleep(waitTime)
	}
	logger.InfoContext(ctx, "[rtc] waiting for the the subsequent nodes to execute STT.")
	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "[rtc] interrupted while waiting for the subsequent nodes to execute STT.")
		return true
	case <-ctx.WaitNodesCancel():
		logger.InfoContext(ctx, "[rtc] the subsequent nodes has finished executing STT.")
		break
	case <-time.After(time.Second * 2):
		logger.WarnContext(ctx, "[rtc] timed out waiting for subsequent nodes to execute STT")
	}
	return false
}

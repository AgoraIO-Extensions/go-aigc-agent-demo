package engine

import (
	"context"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

func (e *Engine) ProcessSendRTC(input <-chan *ttsResult) {
	for {
		r := <-input
		sid, sgid := r.sid, r.sgid
		ctx := r.ctx
		e.sendAudioToRTC(ctx, r.audio, sid, sgid)
	}
}

func (e *Engine) sendAudioToRTC(ctx context.Context, audioChan <-chan []byte, sid, sgid int64) {
	firstSend := true
quickSend:
	for i := 0; i < 18; i++ { // The instantaneous limit for sending packets is 18 packets; otherwise, the average packet rate must be maintained at 1 packet per 10ms.
		var chunk []byte
		var ok bool
		select {
		case chunk, ok = <-audioChan:
			break
		case <-ctx.Done():
			logger.Info("[rtc] Interrupted while sending audio to RTC.", sentencelifecycle.Tag(sid, sgid))
			return
		}
		if !ok {
			logger.Debug("[rtc] Completed sending audio to RTC.", sentencelifecycle.Tag(sid, sgid))
			return
		}
		if firstSend {
			firstSend = false
			sentencelifecycle.SetSidIntoRTC(sid)
			logger.Debug("[rtc] Started sending audio to RTC.", sentencelifecycle.Tag(sid, sgid))
			sentenceGroupBegin := sentencelifecycle.GroupInst().GetAudioEndTime(sid)
			if sentenceGroupBegin == nil {
				logger.Error("Failed to retrieve the start time of the sentence lifecycle group based on SGID.", sentencelifecycle.Tag(sid, sgid))
			} else {
				sentencelifecycle.GroupInst().DeleteAudioEndTime(sid)
				dur := time.Now().Sub(*sentenceGroupBegin)
				logger.Info("[sentence]<duration> STT audio end time ——> Sending the first chunk to RTC", sentencelifecycle.Tag(sid, sgid), slog.Int64("dur", dur.Milliseconds()))
			}
		}

		if err := e.rtc.SendPcm(chunk); err != nil {
			logger.Error("[rtc] Failed to send audio to RTC.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
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
			logger.Info("[rtc] The blocking time is too long; executing quickSend.", sentencelifecycle.Tag(sid, sgid))
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
				logger.Info("[rtc] Interrupted while sending audio to RTC.", sentencelifecycle.Tag(sid, sgid))
				return
			}
			if !ok {
				logger.Info("[rtc] Audio sent to RTC completed", sentencelifecycle.Tag(sid, sgid))
				return
			}

			if dur := time.Since(readStart); dur > time.Millisecond*10 {
				logger.Warn("[rtc] While sending audio to RTC, blocking in audio retrieval for more than 10ms.", slog.Int64("dur", dur.Milliseconds()), sentencelifecycle.Tag(sid, sgid))
			}
			if err := e.rtc.SendPcm(chunk); err != nil {
				logger.Error("[rtc] Failed to send audio to RTC.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
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

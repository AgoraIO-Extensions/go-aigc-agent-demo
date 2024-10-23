package engine

import (
	"context"
	"errors"
	"go-aigc-agent-demo/business/aigcCtx"
	"go-aigc-agent-demo/business/aigcCtx/sentence"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

type sentenceAudio struct {
	ctx   *aigcCtx.AIGCContext
	audio chan *filter.Chunk
}

type sentenceText struct {
	ctx        *aigcCtx.AIGCContext
	finishSend chan struct{}
	sendFailed chan struct{}
	fullText   chan string
}

type sentenceGroupText struct {
	ctx  *aigcCtx.AIGCContext
	text string // the concatenated value of all texts under the current group
}

func (e *Engine) ProcessSTT(input <-chan *filter.Chunk, output chan *sentenceGroupText) {
	var sentenceAudioQueue = make(chan sentenceAudio, 100)
	go e.groupAudio(input, sentenceAudioQueue)

	var sentenceTextQueue = make(chan *sentenceText, 100)
	go func() {
		for {
			sAudio := <-sentenceAudioQueue
			sText := &sentenceText{
				ctx:        sAudio.ctx,
				finishSend: make(chan struct{}, 1),
				sendFailed: make(chan struct{}, 1),
				fullText:   make(chan string, 1),
			}
			sentenceTextQueue <- sText
			e.sendToSTT(sAudio, sText)
		}
	}()

	e.groupText(sentenceTextQueue, output)
}

// groupAudio split the streaming audio into sentence audio and group the sentence audio (sgid).
func (e *Engine) groupAudio(streamAudioQueue <-chan *filter.Chunk, sentenceAudioQueue chan<- sentenceAudio) {
	var (
		cfg           = config.Inst()
		sid, sgid     int64
		prevSMetaData = new(sentence.MetaData)
		sMetaData     = new(sentence.MetaData)
		audioChan     chan *filter.Chunk
	)

	for {
		chunk := <-streamAudioQueue
		if delay := time.Since(chunk.Time).Milliseconds(); delay > 10 {
			logger.Warn("[stt] the time delay in receiving the audio chunk output from the filter >10ms.", slog.Int64("delay", delay), slog.Int64("sid", chunk.Sid))
		}
		sid = chunk.Sid
		if sgid == 0 {
			sgid = sid
		}
		switch chunk.Status {
		case filter.MuteToSpeak:
			*prevSMetaData = *sMetaData
			sMetaData = &sentence.MetaData{Sid: sid}

			ctx := aigcCtx.NewContext(context.WithValue(context.Background(), logger.SentenceMetaData, sMetaData), sMetaData)
			if cfg.InterruptStage == config.AfterFilter {
				ctx.Interrupt()
			}
			sgid = Grouping(prevSMetaData, sid)
			sMetaData.Sgid = sgid
			logger.InfoContext(ctx, "[stt] Get the sentence audio head from upstream")
			audioChan = make(chan *filter.Chunk, 100)
			audioChan <- chunk
			sentenceAudioQueue <- sentenceAudio{
				ctx:   ctx,
				audio: audioChan,
			}
		case filter.SpeakToMute:
			logger.Info("[stt] get the sentence audio tail from upstream", slog.Int64("sid", sid), slog.Int64("sgid", sgid))
			sMetaData.FilterAudioTailRcvTime = chunk.Time
			audioChan <- chunk
			close(audioChan)
		default:
			audioChan <- chunk
		}
	}
}

// sendToSTT send sentence audio to STT and interrupt based on the STT results
func (e *Engine) sendToSTT(sentenceAudio sentenceAudio, sText *sentenceText) {
	ctx := sText.ctx
	sttClient, err := e.sttFactory.CreateSTT(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "[stt] Failed to obtain STT connection instance.", slog.Any("err", err))
		return
	}

	var sendEnd time.Time

	go func() {
		for {
			chunk := <-sentenceAudio.audio
			if chunk.Status == filter.SpeakToMute {
				if dur := time.Since(chunk.Time).Milliseconds(); dur > 10 {
					logger.WarnContext(ctx, "[stt]<duration> The audio chunk took more than 10ms from VAD output to STT input.", slog.Int64("dur", dur))
				}
				sendEnd = time.Now()
				if err = sttClient.Send(nil, true); err != nil {
					sText.sendFailed <- struct{}{}
					logger.ErrorContext(ctx, "[stt] Failed to send the stop command to STT.", slog.Any("err", err))
					return
				}
				sText.finishSend <- struct{}{}
				break
			}
			if err = sttClient.Send(chunk.Data, false); err != nil {
				sText.sendFailed <- struct{}{}
				logger.ErrorContext(ctx, "[stt] Failed to send chunk to STT.", slog.Any("err", err))
				return
			}
		}
	}()

	go func() {
		cfg := config.Inst()
		sttResult := sttClient.GetResult()
		var firstContent string
		for {
			r := <-sttResult
			if r.Fail {
				sText.ctx.ReleaseCtxNode()
				sText.fullText <- ""
				logger.ErrorContext(ctx, "[stt] Asynchronous recognition sText failed")
				return
			}

			if cfg.InterruptStage == config.AfterSTT && firstContent == "" && r.Text != "" {
				sText.ctx.Interrupt()
				logger.InfoContext(ctx, "[stt] do interrupt")
			}

			if firstContent == "" && r.Text != "" {
				firstContent = r.Text
			}

			if r.Complete {
				sText.fullText <- r.Text
				if r.Text == "" {
					if firstContent != "" {
						// It seems that Microsoft’s SDK has encountered this bug before.
						logger.ErrorContext(ctx, "[stt] The STT SDK returned content that was not as expected, it is likely a bug in the SDK")
					}
					sText.ctx.ReleaseCtxNode()
					logger.InfoContext(ctx, "[stt] STT returned an empty string")
					return
				}
				logger.InfoContext(ctx, "[stt]<duration> Received recognized text from STT.", slog.Int64("dur", time.Since(sendEnd).Milliseconds()), slog.String("text", r.Text))
				break
			}
		}
	}()
}

// groupText group sentence texts, and sentences in the same group will be concatenated together
func (e *Engine) groupText(sentenceTextQueue chan *sentenceText, sentenceTextGroupQueue chan *sentenceGroupText) {
	var concatenatedText string
	for {
		/* get one stt recognized text */
		sText := <-sentenceTextQueue
		ctx := sText.ctx
		sMetaData := ctx.MetaData
		if sMetaData.Sid == sMetaData.Sgid { // means it‘s a new group, so reset concatenatedText
			concatenatedText = ""
		}

		/* wait for sAudio to finish sending */
		select {
		case <-sText.sendFailed:
			logger.ErrorContext(ctx, "[stt] sText.sendFailed")
			continue
		case <-sText.finishSend:
			break
		}

		/* concat stt recognized texts that belongs to a group */
		select {
		case <-time.After(time.Second * 5):
			logger.ErrorContext(ctx, "[stt] Timeout waiting for STT to retrieve recognition result: 5 seconds.")
			continue
		case fullText := <-sText.fullText:
			if fullText == "" {
				continue
			}
			concatenatedText = concatenatedText + fullText
			logger.InfoContext(ctx, "[stt] Text after concatenation", slog.String("text", concatenatedText))
		}

		/* check if interrupted */
		if errors.Is(ctx.Err(), context.Canceled) {
			logger.InfoContext(ctx, "[stt] After collecting the STT results, the sentence was interrupted.")
			continue
		}

		sentenceTextGroupQueue <- &sentenceGroupText{
			ctx:  ctx,
			text: concatenatedText,
		}
	}
}

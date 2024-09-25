package engine

import (
	"context"
	"errors"
	"fmt"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/business/interrupt"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"runtime"
	"time"
)

type sentenceAudio struct {
	ctxNode *interrupt.CtxNode
	sid     int64
	sgid    int64
	audio   chan *filter.Chunk
}

type sentenceText struct {
	ctxNode  *interrupt.CtxNode
	sid      int64
	sgid     int64
	fullText chan string
}

type sentenceGroupText struct {
	ctx  context.Context
	sid  int64  // 当前group下最大（最晚）的sid值
	sgid int64  // 当前group下最小（最早）的sid值
	text string // 当前group下所有文本的拼接值
}

func (e *Engine) ProcessSTT(input <-chan *filter.Chunk, output chan *sentenceGroupText) {
	var sentenceAudioQueue = make(chan sentenceAudio, 100)
	go e.groupAudio(input, sentenceAudioQueue)

	var sentenceTextQueue = make(chan *sentenceText, 100)
	go func() {
		for {
			audio := <-sentenceAudioQueue
			st := &sentenceText{
				ctxNode:  audio.ctxNode,
				sid:      audio.sid,
				sgid:     audio.sgid,
				fullText: make(chan string, 1),
			}
			sentenceTextQueue <- st
			e.sendToSTT(audio, st)
		}
	}()

	e.groupText(sentenceTextQueue, output)
}

// groupAudio 将流式音频分开为句子音频，并将将句子音频分组（sgid）
func (e *Engine) groupAudio(streamAudioQueue <-chan *filter.Chunk, sentenceAudioQueue chan<- sentenceAudio) {
	var (
		cfg                 = config.Inst()
		sid                 int64
		sgid                = sentencelifecycle.FirstSid
		audioChan           chan *filter.Chunk
		prevSentenceEndTime = time.Time{}
	)

	for {
		chunk, ok := <-streamAudioQueue
		if !ok {
			logger.Info("[filter] Filter output queue has been consumed and closed.")
			return
		}

		sid = chunk.Sid
		switch chunk.Status {
		case filter.MuteToSpeak:
			ctxNode := interrupt.NewCtxNode(sid)
			if cfg.InterruptStage == config.AfterFilter {
				interrupt.Interrupt(ctxNode)
			}
			sgid = Grouping(sgid, sid, prevSentenceEndTime)
			logger.Info("[stt] Get the sentence audio head from upstream", sentencelifecycle.Tag(sid, sgid))
			sentencelifecycle.GroupInst().SetSidToSgid(sid, sgid)

			audioChan = make(chan *filter.Chunk, 100)
			audioChan <- chunk
			sentenceAudioQueue <- sentenceAudio{
				ctxNode: ctxNode,
				sid:     sid,
				sgid:    sgid,
				audio:   audioChan,
			}
		case filter.SpeakToMute:
			logger.Info("[stt] Get the sentence audio tail from upstream", sentencelifecycle.Tag(sid, sgid))
			prevSentenceEndTime = time.Now()
			audioChan <- chunk
			close(audioChan)
			sentencelifecycle.GroupInst().StoreAudioEndTime(sid, chunk.Time)
		default:
			audioChan <- chunk
		}
	}
}

// sendToSTT 将句子音频发送给stt，并依据stt返回结果进行打断
func (e *Engine) sendToSTT(sentenceAudio sentenceAudio, sText *sentenceText) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Recovered from panic: %v", r))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			logger.Error(fmt.Sprintf("Stack trace:\n%s", buf[:stackSize]))
		}
	}()

	sid := sentenceAudio.sid
	sgid := sentenceAudio.sgid

	sttClient, err := e.sttFactory.CreateSTT(sid)
	if err != nil {
		logger.Error("[stt] Failed to obtain STT connection instance.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
		return
	}

	var sendEnd time.Time

	go func() {
		cfg := config.Inst()
		sttResult := sttClient.GetResult()
		var firstContent string
		for {
			r, ok := <-sttResult
			if !ok {
				logger.Error("[stt] Unreachable code")
				break
			}

			if cfg.InterruptStage == config.AfterSTT && firstContent == "" && r.Text != "" {
				interrupt.Interrupt(sText.ctxNode)
				logger.Info(fmt.Sprintf("[stt] do interrupt, triggered by sid:%d", sid))
			}

			if firstContent == "" && r.Text != "" {
				firstContent = r.Text
			}

			if r.Fail {
				sText.fullText <- ""
				logger.Error("[stt] Asynchronous recognition sText failed", sentencelifecycle.Tag(sid, sgid))
				return
			}
			if r.Complete {
				sText.fullText <- r.Text
				if r.Text == "" {
					if firstContent != "" {
						// It seems that Microsoft’s SDK has encountered this bug before.
						logger.Error("[stt] The STT SDK returned content that was not as expected, it is likely a bug in the SDK")
					}
					if cfg.InterruptStage == config.AfterSTT {
						interrupt.ReleaseCtxNode(sText.ctxNode)
					}
					logger.Info("[stt] STT returned an empty string", sentencelifecycle.Tag(sid, sgid))
					return
				}
				logger.Info("[stt]<duration> Received recognized text from STT.", slog.Int64("dur", time.Since(sendEnd).Milliseconds()), slog.String("text", r.Text), sentencelifecycle.Tag(sid, sgid))
				break
			}
		}
	}()

	for {
		chunk, ok := <-sentenceAudio.audio
		if !ok {
			logger.Error("[stt] Unreachable code")
			break
		}
		if chunk.Status == filter.SpeakToMute {
			if dur := time.Since(chunk.Time).Milliseconds(); dur > 10 {
				logger.Warn("[stt]<duration> The audio chunk took more than 10ms from VAD output to STT input.", slog.Int64("dur", dur), sentencelifecycle.Tag(sid, sgid))
			}
			sendEnd = time.Now()
			if err = sttClient.Send(nil, true); err != nil {
				logger.Error("[stt] Failed to send the stop command to STT.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
				return
			}
			break
		}
		if err = sttClient.Send(chunk.Data, false); err != nil {
			logger.Error("[stt] Failed to send chunk to STT.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
			return
		}
	}
}

// groupText 将句子文本分组，同一组的句子文本会拼接到一起
func (e *Engine) groupText(sentenceTextQueue chan *sentenceText, sentenceTextGroupQueue chan *sentenceGroupText) {
	var concatenatedText string
	for {
		/* get one stt recognized text */
		sText, ok := <-sentenceTextQueue
		if !ok {
			logger.Info("[stt] STT has been closed.")
			return
		}
		sid, sgid := sText.sid, sText.sgid
		if sid == sgid { // means it‘s a new group, so reset concatenatedText
			concatenatedText = ""
		}

		/* concat stt recognized texts that belongs to a group  */
		ctx := sText.ctxNode.Ctx
		select {
		case <-time.After(time.Second * 5):
			logger.Info("[stt] Timeout waiting for STT to retrieve recognition result: 5 seconds.", sentencelifecycle.Tag(sid, sgid))
			continue
		case fullText := <-sText.fullText:
			if fullText == "" {
				continue
			}
			concatenatedText = concatenatedText + fullText
			logger.Info("[stt] Text after concatenation", slog.String("text", concatenatedText), sentencelifecycle.Tag(sid, sgid))
		}

		/* check if interrupted */
		if errors.Is(ctx.Err(), context.Canceled) {
			logger.Info("[stt] After collecting the STT results, the sentence was interrupted.", sentencelifecycle.Tag(sid, sgid))
			continue
		}

		sentenceTextGroupQueue <- &sentenceGroupText{
			ctx:  sText.ctxNode.Ctx,
			sid:  sid,
			sgid: sgid,
			text: concatenatedText,
		}
	}
}

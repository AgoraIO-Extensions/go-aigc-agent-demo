package engine

import (
	"context"
	"errors"
	"fmt"
	"go-aigc-agent-demo/business/exit"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/business/interrupt"
	"go-aigc-agent-demo/business/llm"
	"go-aigc-agent-demo/business/rtc"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/business/stt"
	"go-aigc-agent-demo/business/tts"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"runtime"
	"time"
)

type Engine struct {
	exitWrapper *exit.ExitManager
	filter      *filter.Filter
	rtc         *rtc.RTC
	sttFactory  *stt.Factory
	ttsFactory  *tts.Factory
	llm         *llm.LLM
}

func InitEngine() (*Engine, error) {
	cfg := config.Inst()
	e := &Engine{}

	var err error

	// 初始化「vad」
	e.filter = filter.NewFilter(sentencelifecycle.FirstSid, cfg.Filter.Vad.StartWin, cfg.Filter.Vad.StopWin)

	// 初始化「rtc」
	e.rtc = rtc.NewRTC(cfg.RTC.AppID, "", cfg.RTC.ChannelName, cfg.RTC.UserID, cfg.RTC.Region)
	logger.Info("RTC initialization succeeded")

	// 初始化「stt」
	if e.sttFactory, err = stt.NewFactory(cfg.STT.Select, cfg.STT); err != nil {
		return nil, fmt.Errorf("[stt.NewFactory]%v", err)
	}
	logger.Info("STT initialization succeeded")

	// 初始化 tts
	if e.ttsFactory, err = tts.NewFactory(cfg.TTS.Select, 2); err != nil {
		return nil, fmt.Errorf("初始化tts失败.%v", err)
	}
	logger.Info("TTS initialization succeeded")

	// 初始化 [exit]
	e.exitWrapper = exit.NewExitManager(cfg.StartTime, cfg.MaxLifeTime)

	// 初始化「llm」
	e.llm, err = llm.NewLLM(cfg.LLM.ModelSelect, cfg.LLM.Prompt.Generate(), &cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("[llm.NewLLM]%v", err)
	}

	return e, nil
}

func (e *Engine) Run() error {
	// 异步地：达到最大生命后期后自动退出
	e.exitWrapper.HandlerMaxLifeTime()

	// 注册用户离开事件处理函数
	e.rtc.SetOnUserLeft(e.exitWrapper.OnUserLeft)

	// 注册从rtc接收到的音频处理函数
	e.rtc.SetOnReceiveAudio(e.filter.OnRcvRTCAudio)

	// 处理filter音频

	sttResultQueue := make(chan *STTResult, 20)
	go e.HandlerFilterAudio(sttResultQueue)

	go e.HandleSTTResults(sttResultQueue)

	// rtc建立连接
	if err := e.rtc.Connect(); err != nil {
		return fmt.Errorf("[rtc.Connect]%v", err)
	}
	return nil
}

type STTResult struct {
	ctxNode  *interrupt.CtxNode
	sid      int64
	sgid     int64
	fullText chan string
}

func (e *Engine) HandlerFilterAudio(sttResultQueue chan<- *STTResult) {
	filterAudio := e.filter.OutputAudio()
	var (
		cfg                 = config.Inst()
		sid                 int64
		sgid                = sentencelifecycle.FirstSid
		ctxNode             = interrupt.NewCtxNode(0)
		sentenceAudio       chan *filter.Chunk
		prevSentenceEndTime = time.Time{}
	)
	for {
		chunk, ok := <-filterAudio
		if !ok {
			logger.Info("[filter] Filter output queue has been consumed and closed.")
			return
		}

		sid = chunk.Sid
		switch chunk.Status {
		case filter.MuteToSpeak:
			if cfg.InterruptStage == config.AfterFilter {
				interrupt.Interrupt(ctxNode)
			}
			ctxNode = interrupt.NewCtxNode(sid)
			sgid = Grouping(sgid, sid, prevSentenceEndTime)
			logger.Info("[stt] Get the sentence audio head from upstream", sentencelifecycle.Tag(sid, sgid))
			sentencelifecycle.GroupInst().SetSidToSgid(sid, sgid)
			sentenceSTTResult := &STTResult{ctxNode: ctxNode, sid: sid, sgid: sgid, fullText: make(chan string, 1)}
			sttResultQueue <- sentenceSTTResult

			sentenceAudio = make(chan *filter.Chunk, 100)
			go e.SendAudioToSTT(sid, sgid, sentenceAudio, sentenceSTTResult)
			sentenceAudio <- chunk
		case filter.SpeakToMute:
			logger.Info("[stt] Get the sentence audio tail from upstream", sentencelifecycle.Tag(sid, sgid))
			prevSentenceEndTime = time.Now()
			sentenceAudio <- chunk
			close(sentenceAudio)
			sentencelifecycle.GroupInst().StoreAudioEndTime(sid, chunk.Time)
		default:
			sentenceAudio <- chunk
		}
	}
}

func (e *Engine) SendAudioToSTT(sid, sgid int64, sentenceAudio <-chan *filter.Chunk, sentenceSTTResult *STTResult) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Recovered from panic: %v", r))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			logger.Error(fmt.Sprintf("Stack trace:\n%s", buf[:stackSize]))
		}
	}()

	sttClient, err := e.sttFactory.CreateSTT(sid)
	if err != nil {
		logger.Error("[stt] Failed to obtain STT connection instance.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
		return
	}

	var sendEnd time.Time

	go func() {
		cfg := config.Inst()
		sttResult := sttClient.GetResult()
		var interrupted bool
		var firstContent string
		for {
			r, ok := <-sttResult
			if !ok {
				logger.Error("[stt] Unreachable code")
				break
			}
			if cfg.InterruptStage == config.AfterSTT && !interrupted && r.Text != "" {
				firstContent = r.Text
				interrupt.Interrupt(sentenceSTTResult.ctxNode)
				logger.Info(fmt.Sprintf("[stt] do interrupt, triggered by sid:%d", sid))
				interrupted = true
			}
			if r.Fail {
				sentenceSTTResult.fullText <- ""
				logger.Error("[stt] Asynchronous recognition result failed", sentencelifecycle.Tag(sid, sgid))
				return
			}
			if r.Complete {
				sentenceSTTResult.fullText <- r.Text
				if r.Text == "" {
					if firstContent != "" {
						// It seems that Microsoft’s SDK has encountered this bug before.
						logger.Error("[stt] The STT SDK returned content that was not as expected, it is likely a bug in the SDK")
					}
					if cfg.InterruptStage == config.AfterSTT {
						interrupt.ReleaseCtxNode(sentenceSTTResult.ctxNode)
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
		chunk, ok := <-sentenceAudio
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

func (e *Engine) HandleSTTResults(sttResults <-chan *STTResult) {
	var concatenatedText string
	for {
		/* get one stt recognized text */
		r, ok := <-sttResults
		if !ok {
			logger.Info("[stt] STT has been closed.")
			return
		}
		sid, sgid := r.sid, r.sgid
		if sid == sgid { // means it‘s a new group, so reset concatenatedText
			concatenatedText = ""
		}

		/* concat stt recognized texts that belongs to a group  */
		ctx := r.ctxNode.Ctx
		select {
		case <-time.After(time.Second * 5):
			logger.Info("[stt] Timeout waiting for STT to retrieve recognition result: 5 seconds.", sentencelifecycle.Tag(sid, sgid))
			continue
		case sentenceText := <-r.fullText:
			if sentenceText == "" {
				continue
			}
			concatenatedText = concatenatedText + sentenceText
			logger.Info("[stt] Text after concatenation", slog.String("text", concatenatedText), sentencelifecycle.Tag(sid, sgid))
		}

		/* check if interrupted */
		if errors.Is(ctx.Err(), context.Canceled) {
			logger.Info("[stt] After collecting the STT results, the sentence was interrupted.", sentencelifecycle.Tag(sid, sgid))
			continue
		}

		/* use {$concatenatedText} ask LLM */
		segChan, err := e.llm.Ask(ctx, sid, concatenatedText)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("[llm] Interrupted while requesting LLM.", slog.Any("msg", err), sentencelifecycle.Tag(sid, sgid))
				continue
			}
			logger.Error("[llm.Ask]fail", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
			continue
		}

		/* create a TTS client connection */
		ttsClient, err := e.ttsFactory.CreateTTS(sid)
		if err != nil {
			logger.Error("[tts] Failed to create TTS client instance.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
			continue
		}

		/* asynchronously send the TTS result to RTC */
		go e.SendAudioToRTC(ctx, ttsClient.GetResult(), sid, sgid)

		/* send the LLM result to TTS */
	LOOP:
		for i := 0; ; i++ {
			select {
			case seg, ok := <-segChan:
				if !ok {
					ttsClient.Send(ctx, i, "")
					break LOOP
				}
				ttsClient.Send(ctx, i, seg)
			case <-ctx.Done(): // interrupted
				logger.Info("[tts] The process of sending a segment to TTS was interrupted.", sentencelifecycle.Tag(sid))
				break LOOP
			}
		}
	}
}

func (e *Engine) SendAudioToRTC(ctx context.Context, audioChan <-chan []byte, sid, sgid int64) {
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

package engine

import (
	"context"
	"errors"
	"fmt"
	"go-aigc-agent-demo/business/exit"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/business/llm"
	"go-aigc-agent-demo/business/rtc"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/business/stt"
	"go-aigc-agent-demo/business/tts"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type Engine struct {
	exitWrapper *exit.ExitManager
	filter      *filter.Filter
	rtc         *rtc.RTC
	stt         *stt.STT
	tts         *tts.TTS
	llm         *llm.LLM
}

func InitEngine() (*Engine, error) {
	cfg := config.Inst()
	e := &Engine{}

	var err error

	// 初始化「vad」
	e.filter = filter.NewFilter(sentencelifecycle.FirstSid)

	// 初始化「rtc」
	e.rtc = rtc.NewRTC(cfg.RTC.AppID, "", cfg.RTC.ChannelName, cfg.RTC.UserID, cfg.RTC.Region)
	logger.Inst().Info("初始化rtc成功...")

	// 初始化「stt」
	if e.stt, err = stt.NewSTT(cfg.STT.Select, cfg.STT); err != nil {
		return nil, fmt.Errorf("[stt.NewSTT]%v", err)
	}
	logger.Inst().Info("初始化stt成功")

	// 初始化 tts
	if e.tts, err = tts.Init(); err != nil {
		return nil, fmt.Errorf("初始化tts失败.%v", err)
	}
	logger.Inst().Info("初始化tts成功...")

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
	go func() {
		e.HandlerFilterAudio()
	}()

	// rtc建立连接
	if err := e.rtc.Connect(); err != nil {
		return fmt.Errorf("[rtc.Connect]%v", err)
	}
	return nil
}

type STTResult struct {
	ctx      context.Context
	sid      int64
	sgid     int64
	fullText chan string
}

func (e *Engine) HandlerFilterAudio() {
	filterAudio := e.filter.OutputAudio()
	var (
		sid           int64
		sgid          = sentencelifecycle.FirstSid
		ctx, cancel   = context.WithCancel(context.Background())
		sentenceAudio chan *filter.Chunk
		sttResults    = make(chan *STTResult, 20)
	)
	go e.HandleSTTResults(sttResults)
	for {
		chunk, ok := <-filterAudio
		if !ok {
			cancel()
			logger.Inst().Info("[filter] filter输出队列已消费完毕并关闭")
			return
		}

		sid = chunk.Sid
		switch chunk.Status {
		case filter.MuteToSpeak:
			cancel() // 触发断句，打断所有历史sentence的处理逻辑
			ctx, cancel = context.WithCancel(context.Background())
			if sentencelifecycle.IfSidIntoRTC(sid - 1) { // 判断上一个sid对应的sentence是否已经到达发送rtc的阶段，如果到了，则开启新的 sentencelifecycle group
				sgid = sid // 新的group
				sentencelifecycle.DeleteSidIntoRtc(sid - 1)
			}
			logger.Inst().Info("[stt] 从上游取出sentence音频头", sentencelifecycle.Tag(sid, sgid))
			sentencelifecycle.GroupInst().SetSidToSgid(sid, sgid)
			sentenceSTTResult := &STTResult{ctx: ctx, sid: sid, sgid: sgid, fullText: make(chan string, 1)}
			sttResults <- sentenceSTTResult

			sentenceAudio = make(chan *filter.Chunk, 100)
			go e.HandlerSentenceAudio(sid, sgid, sentenceAudio, sentenceSTTResult)
			sentenceAudio <- chunk
		case filter.SpeakToMute:
			sentenceAudio <- chunk
			close(sentenceAudio)
			sentencelifecycle.GroupInst().StoreInAudioEndTimeInOneSentenceGroup(sgid, chunk.Time)
		default:
			sentenceAudio <- chunk
		}
	}
}

func (e *Engine) HandlerSentenceAudio(sid, sgid int64, sentenceAudio <-chan *filter.Chunk, sentenceSTTResult *STTResult) {
	defer func() {
		if r := recover(); r != nil {
			logger.Inst().Error(fmt.Sprintf("Recovered from panic: %v", r))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			logger.Inst().Error(fmt.Sprintf("Stack trace:\n%s", buf[:stackSize]))
		}
	}()

	sttConn, err := e.stt.GetOneConnection(sid)
	if err != nil {
		logger.Inst().Error("[stt] 获取stt连接实例失败", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
		return
	}

	var sendEnd time.Time

	go func() {
		sttResult := sttConn.GetResult()
		for {
			r, ok := <-sttResult
			if !ok {
				logger.Inst().Error("[stt] 理论不可达代码")
				break
			}
			if r.Fail {
				sentenceSTTResult.fullText <- ""
				logger.Inst().Error("[stt] 异步识别结果失败", sentencelifecycle.Tag(sid, sgid))
				return
			}
			if r.Complete {
				sentenceSTTResult.fullText <- r.Text
				if r.Text == "" {
					logger.Inst().Info("[stt] stt返回空字符串", sentencelifecycle.Tag(sid, sgid))
					return
				}
				logger.Inst().Info("[stt]<duration> 从stt收到识别文本", zap.Int64("dur", time.Since(sendEnd).Milliseconds()), zap.String("text", r.Text), sentencelifecycle.Tag(sid, sgid))
				break
			}
		}
	}()

	for {
		chunk, ok := <-sentenceAudio
		if !ok {
			logger.Inst().Error("[stt] 理论不可达代码")
			break
		}
		if chunk.Status == filter.SpeakToMute {
			if dur := time.Since(chunk.Time).Milliseconds(); dur > 10 {
				logger.Inst().Warn("[stt]<duration> 音频chunk在 vad输出 ——> stt输入 的过程耗时>10ms", zap.Int64("dur", dur), sentencelifecycle.Tag(sid, sgid))
			}
			sendEnd = time.Now()
			if err = sttConn.Send(nil, true); err != nil {
				logger.Inst().Error("[stt] 往stt发送 stop指令 标记失败", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
				return
			}
			break
		}
		if err = sttConn.Send(chunk.Data, false); err != nil {
			logger.Inst().Error("[stt] 往stt发送chunk失败", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
			return
		}
	}
}

func (e *Engine) HandleSTTResults(sttResults <-chan *STTResult) {
	var groupSentencesMerge string
	for {
		r, ok := <-sttResults
		if !ok {
			logger.Inst().Info("[stt] stt已关闭")
			return
		}
		sid, sgid := r.sid, r.sgid
		if sid == sgid { // 新的group
			groupSentencesMerge = ""
		}

		ctx := r.ctx
		select {
		case <-time.After(time.Second * 5):
			logger.Inst().Info("[stt] 等待stt获取识别结果超时5s", sentencelifecycle.Tag(sid, sgid))
			continue
		case sentenceText := <-r.fullText:
			groupSentencesMerge = groupSentencesMerge + sentenceText
			if groupSentencesMerge == "" {
				logger.Inst().Info("[stt] 合并后的stt结果为空字符串", sentencelifecycle.Tag(sid, sgid))
				continue
			}
			logger.Inst().Info("[stt] 合并后的文本", zap.String("text", groupSentencesMerge), sentencelifecycle.Tag(sid, sgid))
		}

		if errors.Is(ctx.Err(), context.Canceled) {
			logger.Inst().Info("[stt] 合并文本后，被打断", sentencelifecycle.Tag(sid, sgid))
			continue
		}

		segChan, err := e.llm.Ask(ctx, sid, groupSentencesMerge)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Inst().Info("[llm] 请求llm时被打断", zap.String("msg", err.Error()), sentencelifecycle.Tag(sid, sgid))
				continue
			}
			logger.Inst().Error("[llm.Ask]fail", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
			continue
		}
		ttsSender := e.tts.NewSender(sid)

		go e.SendAudioToRTC(ctx, ttsSender.GetResult(), sid, sgid)
	LOOP:
		for i := 0; ; i++ {
			select {
			case seg, ok := <-segChan:
				if !ok {
					ttsSender.Send(ctx, i, "")
					break LOOP
				}
				ttsSender.Send(ctx, i, seg)
			case <-ctx.Done():
				logger.Inst().Info("[tts] 往tts发segment的过程被打断", sentencelifecycle.Tag(sid))
				break LOOP
			}
		}
	}
}

func (e *Engine) SendAudioToRTC(ctx context.Context, audioChan <-chan []byte, sid, sgid int64) {
	firstSend := true
quickSend:
	for i := 0; i < 18; i++ { // 瞬时发送包的个数上限是18个包，除此之外，其他时间段内平均发包速率要保持在1个包/10ms
		var chunk []byte
		var ok bool
		select {
		case chunk, ok = <-audioChan:
			break
		case <-ctx.Done():
			logger.Inst().Info("[rtc] 在往rtc发送音频时被打断", sentencelifecycle.Tag(sid, sgid))
			return
		}
		if !ok {
			logger.Inst().Debug("[rtc] 往rtc发送 音频 完毕", sentencelifecycle.Tag(sid, sgid))
			return
		}
		if firstSend {
			firstSend = false
			sentencelifecycle.SetSidIntoRTC(sid)
			logger.Inst().Debug("[rtc] 开始往rtc发送音频", sentencelifecycle.Tag(sid, sgid))
			sentenceGroupBegin := sentencelifecycle.GroupInst().GetInAudioEndTimeInOneSentenceGroup(sgid)
			if sentenceGroupBegin == nil {
				logger.Inst().Error("根据sgid获取 sentencelifecycle group 开始时间失败", sentencelifecycle.Tag(sid, sgid))
			} else {
				sentencelifecycle.GroupInst().DeleteInAudioEndTimeInOneSentenceGroup(sgid)
				dur := time.Now().Sub(*sentenceGroupBegin)
				logger.Inst().Info("[sentence]<duration> stt发送音频结束时刻 ——> 往rtc发送第一个chunk", sentencelifecycle.Tag(sid, sgid), zap.Int64("dur", dur.Milliseconds()))
			}
		}

		if err := e.rtc.SendPcm(chunk); err != nil {
			logger.Inst().Error("[rtc] 往rtc发送音频报错", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
			return
		}
	}

	firstSendTime := time.Now()
	sendCount := 0
	shouldSendCount := 0
	for {
		time.Sleep(time.Millisecond * 50)
		shouldSendCount = int(time.Since(firstSendTime).Milliseconds())/10 - sendCount
		if shouldSendCount > 18 { // 如果下面的（<-audioChan）操作阻塞了很久（>=140ms），那么shouldSendCount会>18：(140+50)/10=19>18
			logger.Inst().Info("[rtc] 阻塞时间过长，将重新进行快速发送", sentencelifecycle.Tag(sid, sgid))
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
				logger.Inst().Info("[rtc] 在往rtc发送音频时被打断", sentencelifecycle.Tag(sid, sgid))
				return
			}
			if !ok {
				logger.Inst().Info("[rtc] 往rtc发送 音频 完毕", sentencelifecycle.Tag(sid, sgid))
				return
			}

			if dur := time.Since(readStart); dur > time.Millisecond*10 {
				logger.Inst().Warn("[rtc] 往rtc发送音频时，阻塞在获取音频的逻辑>10ms", zap.Int64("dur", dur.Milliseconds()), sentencelifecycle.Tag(sid, sgid))
			}
			if err := e.rtc.SendPcm(chunk); err != nil {
				logger.Inst().Error("[rtc] 往rtc发送音频报错", zap.Error(err), sentencelifecycle.Tag(sid, sgid))
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

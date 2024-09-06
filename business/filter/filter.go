package filter

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

type Chunk struct {
	Sid    int64
	Data   []byte // 长度是320字节
	Status ResultCode
	Time   time.Time
}

type Filter struct {
	sid    int64
	vad    *Vad
	output chan *Chunk
}

func NewFilter(FirstSid int64) *Filter {
	return &Filter{
		sid:    FirstSid - 1,
		vad:    NewVad(),
		output: make(chan *Chunk, 1000),
	}
}

func (f *Filter) OnRcvRTCAudio(con *agoraservice.RtcConnection, channelId string, uid string, inFrame *agoraservice.PcmAudioFrame) {
	cks, status, err := f.vad.ProcessPcmFrame(inFrame)
	if err != nil {
		logger.Info("[vad] 处理音频失败", slog.Any("err", err))
		return
	}

	now := time.Now()
	switch status {
	case Mute:
		return
	case MuteToSpeak:
		f.sid++
		logger.Info("[filter] 收到sentence音频头", slog.Int64("sid", f.sid))
		for i, ck := range cks {
			state := MuteToSpeak
			if i != 0 {
				state = Speaking
			}
			f.output <- &Chunk{
				Sid:    f.sid,
				Data:   ck,
				Status: state,
				Time:   now,
			}
		}
	case Speaking:
		for _, ck := range cks {
			f.output <- &Chunk{
				Sid:    f.sid,
				Data:   ck,
				Status: Speaking,
				Time:   now,
			}
		}
	case SpeakToMute:
		logger.Info("[filter] 收到sentence音频尾", slog.Int64("sid", f.sid))
		f.output <- &Chunk{
			Sid:    f.sid,
			Status: SpeakToMute,
			Time:   now,
		}
	default:
		logger.Error("[filter] This code should never be executed", slog.Any("status", status))
	}
	return
}

func (f *Filter) OutputAudio() <-chan *Chunk {
	return f.output
}

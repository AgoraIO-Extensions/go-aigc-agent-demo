package filter

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

type Chunk struct {
	Sid    int64
	Data   []byte // The length is 320 bytes.
	Status ResultCode
	Time   time.Time
}

type Filter struct {
	sid    int64
	vad    *Vad
	output chan *Chunk
}

func NewFilter(FirstSid int64, startWin, StopWin int) *Filter {
	return &Filter{
		sid:    FirstSid - 1,
		vad:    NewVad(startWin, StopWin),
		output: make(chan *Chunk, 1000),
	}
}

func (f *Filter) OnRcvRTCAudio(con *agoraservice.RtcConnection, channelId string, uid string, inFrame *agoraservice.PcmAudioFrame) {
	cks, status, err := f.vad.ProcessPcmFrame(inFrame)
	if err != nil {
		logger.Info("[vad] Failed to process audio.", slog.Any("err", err))
		return
	}

	now := time.Now()
	switch status {
	case Mute:
		return
	case MuteToSpeak:
		f.sid++
		logger.Info("[filter] Received sentence audio header.", slog.Int64("sid", f.sid))
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
		logger.Info("[filter] Received sentence audio tail.", slog.Int64("sid", f.sid))
		f.output <- &Chunk{
			Sid:    f.sid,
			Status: SpeakToMute,
			Time:   now,
		}
	default:
		logger.Error("[filter] Unreachable code", slog.Any("status", status))
	}
	return
}

func (f *Filter) OutputAudio() <-chan *Chunk {
	return f.output
}

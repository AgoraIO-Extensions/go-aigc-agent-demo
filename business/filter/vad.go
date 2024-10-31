package filter

import (
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
)

type Vad struct {
	v *agoraservice.AudioVadV2
}

func NewVad(startWin, StopWin int) *Vad {
	cfg := &agoraservice.AudioVadConfigV2{
		PreStartRecognizeCount: 16,
		StartRecognizeCount:    startWin,
		StopRecognizeCount:     StopWin,
		ActivePercent:          0.7,
		InactivePercent:        0.5,
		StartVoiceProb:         70,
		StartRms:               -50.0,
		StopVoiceProb:          70,
		StopRms:                -50.0,
	}
	v := agoraservice.NewAudioVadV2(cfg)
	return &Vad{v: v}
}

func (v *Vad) Release() {
	v.v.Release()
}

type ResultCode int

const (
	Mute        ResultCode = 0
	MuteToSpeak ResultCode = 1
	Speaking    ResultCode = 2
	SpeakToMute ResultCode = 3
)

func retToCode(ret agoraservice.VadState) ResultCode {
	switch ret {
	case agoraservice.VadStateWaitSpeeking:
		return Mute
	case agoraservice.VadStateStartSpeeking:
		return MuteToSpeak
	case agoraservice.VadStateIsSpeeking:
		return Speaking
	case agoraservice.VadStateStopSpeeking:
		return SpeakToMute
	default:
		return ResultCode(ret)
	}
}

func (v *Vad) ProcessPcmFrame(inFrame *agoraservice.AudioFrame) ([][]byte, ResultCode, error) {
	outFrame, ret := v.v.Process(inFrame)
	code := retToCode(ret)
	switch code {
	case Mute, SpeakToMute:
		return nil, code, nil
	case MuteToSpeak, Speaking:
		break
	default:
		return nil, code, fmt.Errorf("unexpected code:%d", code)
	}

	n := len(outFrame.Buffer)
	if n == 0 || n%320 != 0 {
		logger.Warn(fmt.Sprintf("if n=len(outFrame.Data), then n=%d, n%%320=%d; it's unexpected, code:%d", n, n%320, code))
		//return nil, code, fmt.Errorf("if n=len(outFrame.Data), then n=%d, n%%320=%d; it's unexpected, code:%d", n, n%320, code)
	}

	ckNums := n / 320

	if ckNums == 0 {
		return [][]byte{make([]byte, 320)}, code, nil
	}
	var chunks [][]byte
	for i := 0; i < ckNums; i++ {
		ck := make([]byte, 320)
		copy(ck, outFrame.Buffer[i*320:(i+1)*320])
		chunks = append(chunks, ck)
	}

	return chunks, code, nil
}

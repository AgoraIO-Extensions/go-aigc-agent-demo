package filter

import (
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
)

type Vad struct {
	v *agoraservice.AudioVad
}

func NewVad(startWin, StopWin int) *Vad {
	cfg := &agoraservice.AudioVadConfig{
		FftSz:                 1024,
		AnaWindowSz:           768,
		HopSz:                 160,
		FrqInputAvailableFlag: 0,
		UseCVersionAIModule:   0,
		VoiceProbThr:          0.7,
		RmsThr:                -40.0,
		JointThr:              0.0,
		Aggressive:            2.0,

		StartRecognizeCount:    startWin,
		StopRecognizeCount:     StopWin,
		PreStartRecognizeCount: 10,
		ActivePercent:          0.6,
		InactivePercent:        0.2,
	}
	v := agoraservice.NewAudioVad(cfg)
	return &Vad{v: v}
}

func (v *Vad) Release() {
	v.v.Release()
}

type ResultCode int

const (
	Err         ResultCode = -1
	Mute        ResultCode = 0
	MuteToSpeak ResultCode = 1
	Speaking    ResultCode = 2
	SpeakToMute ResultCode = 3
)

func (v *Vad) ProcessPcmFrame(inFrame *agoraservice.PcmAudioFrame) ([][]byte, ResultCode, error) {
	outFrame, ret := v.v.ProcessPcmFrame(inFrame)
	if ret == -1 {
		return nil, Err, fmt.Errorf("code:%d", ret)
	}

	code := ResultCode(ret)
	switch code {
	case Mute, SpeakToMute:
		return nil, code, nil
	case MuteToSpeak, Speaking:
		break
	default:
		return nil, code, fmt.Errorf("unexpected code:%d", code)
	}

	n := len(outFrame.Data)
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
		copy(ck, outFrame.Data[i*320:(i+1)*320])
		chunks = append(chunks, ck)
	}

	return chunks, code, nil
}

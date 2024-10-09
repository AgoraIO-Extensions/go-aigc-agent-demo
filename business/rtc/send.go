package rtc

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"time"
)

func (r *RTC) SendPcm(chunk []byte) error {
	if len(chunk) != 320 {
		return fmt.Errorf("len(chunk) != 320")
	}

	if err := r.sendLimiter.Wait(context.Background()); err != nil {
		return fmt.Errorf("[sendLimiter.Wait]%w", err)
	}

	frame := &agoraservice.PcmAudioFrame{
		Data:              chunk,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}
	if code := r.pcmSender.SendPcmData(frame); code != 0 {
		return fmt.Errorf("err code=%d", code)
	}
	return nil
}

func (r *RTC) SendStreamMessage(msg []byte) {
	r.conn.SendStreamMessage(r.streamID, msg)
	time.Sleep(time.Millisecond)
}

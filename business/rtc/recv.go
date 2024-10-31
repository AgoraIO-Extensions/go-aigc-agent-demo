package rtc

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
)

type OnReceiveAudio func(localUser *agoraservice.LocalUser, channelId string, userId string, frame *agoraservice.AudioFrame) bool

var audioObserver = &agoraservice.AudioFrameObserver{
	OnPlaybackAudioFrameBeforeMixing: nil,
}

func (r *RTC) SetOnReceiveAudio(handlerFunc OnReceiveAudio) {
	audioObserver.OnPlaybackAudioFrameBeforeMixing = handlerFunc
}

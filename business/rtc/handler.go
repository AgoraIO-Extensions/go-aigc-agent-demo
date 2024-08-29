package rtc

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
)

type OnConnected func(conn *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int)
type OnDisconnected func(conn *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int)
type OnUserJoined func(conn *agoraservice.RtcConnection, uid string)
type OnUserLeft func(conn *agoraservice.RtcConnection, uid string, reason int)

func (r *RTC) SetOnConnected(handlerFunc OnConnected) {
	r.connConfig.ConnectionHandler.OnConnected = handlerFunc
}

func (r *RTC) SetOnDisconnected(handlerFunc OnDisconnected) {
	r.connConfig.ConnectionHandler.OnDisconnected = handlerFunc
}

func (r *RTC) SetOnUserJoined(handlerFunc OnUserJoined) {
	r.connConfig.ConnectionHandler.OnUserJoined = handlerFunc
}

func (r *RTC) SetOnUserLeft(handlerFunc OnUserLeft) {
	r.connConfig.ConnectionHandler.OnUserLeft = handlerFunc
}

type OnReceiveAudio func(conn *agoraservice.RtcConnection, channelId string, uid string, frame *agoraservice.PcmAudioFrame)

func (r *RTC) SetOnReceiveAudio(handlerFunc OnReceiveAudio) {
	r.connConfig.AudioFrameObserver = &agoraservice.RtcConnectionAudioFrameObserver{OnPlaybackAudioFrameBeforeMixing: handlerFunc}
}

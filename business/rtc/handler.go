package rtc

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
)

var conHandler = &agoraservice.RtcConnectionObserver{
	OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Connected, reason %d\n", reason)
	},
	OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Disconnected, reason %d\n", reason)
	},
	OnConnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Connecting, reason %d\n", reason)
	},
	OnReconnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Reconnecting, reason %d\n", reason)
	},
	OnReconnected: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Reconnected, reason %d\n", reason)
	},
	OnConnectionLost: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo) {
		logger.Info("[rtc] Connection lost\n")
	},
	OnConnectionFailure: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
		logger.Info("[rtc] Connection failure, error code %d\n", errCode)
	},
	OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
		logger.Info("[rtc] user joined, " + uid)
	},
	OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
		logger.Info("[rtc] user left, " + uid)
	},
}

type OnConnected func(conn *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int)
type OnDisconnected func(conn *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int)
type OnUserJoined func(conn *agoraservice.RtcConnection, uid string)
type OnUserLeft func(conn *agoraservice.RtcConnection, uid string, reason int)

func (r *RTC) SetOnConnected(handlerFunc OnConnected) {
	conHandler.OnConnected = handlerFunc
}

func (r *RTC) SetOnDisconnected(handlerFunc OnDisconnected) {
	conHandler.OnDisconnected = handlerFunc
}

func (r *RTC) SetOnUserJoined(handlerFunc OnUserJoined) {
	conHandler.OnUserJoined = handlerFunc
}

func (r *RTC) SetOnUserLeft(handlerFunc OnUserLeft) {
	conHandler.OnUserLeft = handlerFunc
}

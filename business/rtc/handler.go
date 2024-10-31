package rtc

import (
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

var conHandler = &agoraservice.RtcConnectionObserver{
	OnConnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Connected", slog.Int("reason", reason))
	},
	OnDisconnected: func(con *agoraservice.RtcConnection, info *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Disconnected", slog.Int("reason", reason))
	},
	OnConnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Connecting", slog.Int("reason", reason))
	},
	OnReconnecting: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Reconnecting", slog.Int("reason", reason))
	},
	OnReconnected: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, reason int) {
		logger.Info("[rtc] Reconnected", slog.Int("reason", reason))
	},
	OnConnectionLost: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo) {
		logger.Error("[rtc] Connection lost")
	},
	OnConnectionFailure: func(con *agoraservice.RtcConnection, conInfo *agoraservice.RtcConnectionInfo, errCode int) {
		logger.Error("[rtc] Connection failure", slog.Int("errCode", errCode))
	},
	OnUserJoined: func(con *agoraservice.RtcConnection, uid string) {
		logger.Info("[rtc] user joined", slog.String("uid", uid))
	},
	OnUserLeft: func(con *agoraservice.RtcConnection, uid string, reason int) {
		logger.Info("[rtc] user left", slog.String("uid", uid), slog.Int("reason", reason))
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

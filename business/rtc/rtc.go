package rtc

import (
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"golang.org/x/time/rate"
	"strings"
)

var regionMap = map[string]uint{
	"cn": 0x00000001, // Mainland China
	"na": 0x00000002, // North America
	"eu": 0x00000004, // Europe
	//"as":   0x00000008, // Asia (excluding Mainland China)
	"ap":   0x00000008, // Asia-Pacific
	"jp":   0x00000010, // Japan
	"in":   0x00000020, // India
	"glob": 0xFFFFFFFF, // Global (default value)
}

type InitParams struct {
	appid       string
	token       string
	channelName string
	userID      string
	AreaCode    uint
}

type RTC struct {
	initParams  *InitParams
	connConfig  *agoraservice.RtcConnectionConfig
	conn        *agoraservice.RtcConnection
	pcmSender   *agoraservice.PcmSender
	sendLimiter *rate.Limiter
	streamID    int
}

func NewRTC(appid, token, channelName, userId, region string) *RTC {
	region = strings.ToLower(region)
	areaCode := regionMap["glob"]
	if code, ok := regionMap[region]; ok {
		areaCode = code
	}

	svcCfg := agoraservice.AgoraServiceConfig{
		AppId:         appid,
		AudioScenario: agoraservice.AUDIO_SCENARIO_CHORUS,
		LogPath:       "./agora_rtc_log/agorasdk.log",
		LogSize:       512 * 1024,
		AreaCode:      areaCode,
	}
	agoraservice.Init(&svcCfg)

	connCfg := &agoraservice.RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,

		SubAudioConfig: &agoraservice.SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
		ConnectionHandler:  &agoraservice.RtcConnectionEventHandler{},
		AudioFrameObserver: nil,
		VideoFrameObserver: nil,
	}

	params := &InitParams{
		appid:       appid,
		token:       token,
		channelName: channelName,
		userID:      userId,
		AreaCode:    areaCode,
	}

	return &RTC{
		initParams:  params,
		connConfig:  connCfg,
		sendLimiter: rate.NewLimiter(100, 18),
	}
}

func (r *RTC) Connect() error {
	r.conn = agoraservice.NewConnection(r.connConfig)
	r.pcmSender = r.conn.NewPcmSender()
	code := r.conn.Connect(r.initParams.token, r.initParams.channelName, r.initParams.userID)
	if code != 0 {
		return fmt.Errorf("err code:%d", code)
	}
	r.pcmSender.Start()
	r.pcmSender.AdjustVolume(100)
	r.streamID, code = r.conn.CreateDataStream(true, true)
	if code != 0 {
		return fmt.Errorf("[CreateDataStream] err code:%d", code)
	}
	return nil
}

func (r *RTC) Close() {
	r.pcmSender.Stop()
	r.conn.Disconnect()
	agoraservice.Destroy()
	r.pcmSender.Release()
	r.conn.Release()
}

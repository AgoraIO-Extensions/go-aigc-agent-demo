package rtc

import (
	"fmt"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"strings"
)

var regionMap = map[string]uint{
	"cn": 0x00000001, // 中国大陆
	"na": 0x00000002, // 北美
	"eu": 0x00000004, // 欧洲
	//"as":   0x00000008, // 亚洲（不包括中国大陆）
	"ap":   0x00000008, // 亚太
	"jp":   0x00000010, // 日本
	"in":   0x00000020, // 印度
	"glob": 0xFFFFFFFF, // 全球（默认值）
}

type InitParams struct {
	appid       string
	token       string
	channelName string
	userID      string
	Region      uint
}

type RTC struct {
	initParams *InitParams
	connConfig *agoraservice.RtcConnectionConfig
	conn       *agoraservice.RtcConnection
	pcmSender  *agoraservice.PcmSender
	streamID   int
}

func NewRTC(appid, token, channelName, userId, region string) *RTC {
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

	region = strings.ToLower(region)
	if _, ok := regionMap[region]; !ok {
		region = "glob"
	}
	params := &InitParams{
		appid:       appid,
		token:       token,
		channelName: channelName,
		userID:      userId,
		Region:      regionMap[region],
	}

	return &RTC{initParams: params, connConfig: connCfg}
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

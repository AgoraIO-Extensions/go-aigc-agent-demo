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
	initParams       *InitParams
	connConfig       *agoraservice.RtcConnectionConfig
	track            *agoraservice.LocalAudioTrack
	conn             *agoraservice.RtcConnection
	mediaNodeFactory *agoraservice.MediaNodeFactory
	pcmSender        *agoraservice.AudioPcmDataSender
	sendLimiter      *rate.Limiter
	streamID         int
}

func NewRTC(appid, token, channelName, userId, region string) *RTC {
	region = strings.ToLower(region)
	areaCode := regionMap["glob"]
	if code, ok := regionMap[region]; ok {
		areaCode = code
	}

	svcCfg := agoraservice.NewAgoraServiceConfig()
	svcCfg.AppId = appid
	agoraservice.Initialize(svcCfg)

	connCfg := &agoraservice.RtcConnectionConfig{
		AutoSubscribeAudio: true,
		AutoSubscribeVideo: false,
		ClientRole:         agoraservice.ClientRoleBroadcaster,
		ChannelProfile:     agoraservice.ChannelProfileLiveBroadcasting,
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
	r.conn = agoraservice.NewRtcConnection(r.connConfig)
	localUser := r.conn.GetLocalUser()
	localUser.SetPlaybackAudioFrameBeforeMixingParameters(1, 16000)
	localUser.RegisterAudioFrameObserver(audioObserver)

	r.conn.RegisterObserver(conHandler)

	r.mediaNodeFactory = agoraservice.NewMediaNodeFactory()
	r.pcmSender = r.mediaNodeFactory.NewAudioPcmDataSender()
	agoraservice.EnableExtension("agora.builtin", "agora_audio_label_generator", "", true)
	agoraservice.GetAgoraParameter().SetParameters("{\"che.audio.label.enable\": true}")
	r.track = agoraservice.NewCustomAudioTrackPcm(r.pcmSender)
	localUser.SetAudioScenario(agoraservice.AudioScenarioChorus)
	code := r.conn.Connect(r.initParams.token, r.initParams.channelName, r.initParams.userID)
	if code != 0 {
		return fmt.Errorf("err code:%d", code)
	}
	r.track.SetEnabled(true)
	localUser.PublishAudio(r.track)
	r.track.AdjustPublishVolume(100)
	r.streamID, code = r.conn.CreateDataStream(true, true)
	if code != 0 {
		return fmt.Errorf("[CreateDataStream] err code:%d", code)
	}
	return nil
}

func (r *RTC) Release() {
	r.track.Release()
	r.pcmSender.Release()
	r.conn.Release()
	r.mediaNodeFactory.Release()
	agoraservice.Release()
}

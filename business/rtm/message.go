package rtm

import (
	"fmt"
	"go-aigc-agent-demo/business/aigcCtx"
	"go-aigc-agent-demo/business/rtc"
	aigcmsg "go-aigc-agent-demo/business/rtm/proto"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

var rtmSend *RtmSend

type RtmSend struct {
	version     int32
	userID      int32
	baseRoundID int64
	rtc         *rtc.RTC
}

func Init(ver, uid int32, firstSID int64, rtc *rtc.RTC) {
	logger.Info(fmt.Sprintf("[message] baseRoundID is %d", firstSID))
	rtmSend = &RtmSend{
		version:     ver,
		userID:      uid,
		baseRoundID: firstSID,
		rtc:         rtc,
	}
}

func (r *RtmSend) buildMsg(tp int32, ver, userID, flag int32, sid int64, content string) ([]byte, error) {
	msg := aigcmsg.AigcMessage{
		Magicnum:  1,
		Type:      tp,
		Version:   ver,
		Userid:    userID,
		Roundid:   int32(sid - r.baseRoundID),
		Flag:      flag,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	}
	return proto.Marshal(&msg)
}

// buildSTTMsg
// flag: 0表示正在返回，此时content有内容；1表示返回结束，此时content为空
func (r *RtmSend) buildSTTMsg(flag int32, sid int64, content string) ([]byte, error) {
	return r.buildMsg(100, r.version, r.userID, flag, sid, content)
}

// buildLLMMsg
// flag: 0表示进行中，此时content有内容；1表示结束，此时content为空
func (r *RtmSend) buildLLMMsg(flag int32, sid int64, content string) ([]byte, error) {
	return r.buildMsg(120, r.version, r.userID, flag, sid, content)
}

// buildTTSMsg
// flag: 0表示 1个round 内开始推送瞬间；1表示 1个round 推送结束瞬间
func (r *RtmSend) buildTTSMsg(flag int32, sid int64) ([]byte, error) {
	return r.buildMsg(130, r.version, r.userID, flag, sid, "")
}

// buildSentenceLifecycleMsg
// flag: 0表示 round 开始的时刻；1表示 round round的结束瞬间
func (r *RtmSend) buildSentenceLifecycleMsg(flag int32, sid int64) ([]byte, error) {
	return r.buildMsg(140, r.version, r.userID, flag, sid, "")
}

type FlagType int32

const (
	FlagNoFin FlagType = 0
	FlagFin   FlagType = 1
)

func SendSttMsg(ctx *aigcCtx.AIGCContext, flag FlagType, content string) error {
	if !config.Inst().RTC.OpenMsgReturn {
		return nil
	}

	msgBytes, err := rtmSend.buildSTTMsg(int32(flag), ctx.MetaData.Sgid, strings.TrimSpace(content))
	if err != nil {
		return fmt.Errorf("[buildSTTMsg]%v", err)
	}
	rtmSend.rtc.SendStreamMessage(msgBytes)
	logger.InfoContext(ctx, "[rtm] 发送stt消息到rtm成功", zap.Any("flag", flag))
	return nil
}

func SendLlmMsg(ctx *aigcCtx.AIGCContext, flag FlagType, content string) error {
	if !config.Inst().RTC.OpenMsgReturn {
		return nil
	}
	msgBytes, err := rtmSend.buildLLMMsg(int32(flag), ctx.MetaData.Sgid, strings.TrimSpace(content))
	if err != nil {
		return fmt.Errorf("[buildLLMMsg]%v", err)
	}
	rtmSend.rtc.SendStreamMessage(msgBytes)
	logger.InfoContext(ctx, "[rtm] 发送llm消息到rtm成功", zap.Any("flag", flag))
	return nil
}

func SendTtsMsg(ctx *aigcCtx.AIGCContext, flag FlagType) error {
	if !config.Inst().RTC.OpenMsgReturn {
		return nil
	}
	msgBytes, err := rtmSend.buildTTSMsg(int32(flag), ctx.MetaData.Sgid)
	if err != nil {
		return fmt.Errorf("[buildTTSMsg]%v", err)
	}
	rtmSend.rtc.SendStreamMessage(msgBytes)
	logger.InfoContext(ctx, "[rtm] 发送tts消息到rtm成功", zap.Any("flag", flag))
	if flag == 1 {
		if err = SendSentenceMsg(ctx, FlagFin); err != nil {
			return fmt.Errorf("[e.SendSentenceMsg]%v", err)
		}
	}
	return nil
}

func SendSentenceMsg(ctx *aigcCtx.AIGCContext, flag FlagType) error {
	if !config.Inst().RTC.OpenMsgReturn {
		return nil
	}
	sgid := ctx.MetaData.Sgid
	msgBytes, err := rtmSend.buildSentenceLifecycleMsg(int32(flag), sgid)
	if err != nil {
		return fmt.Errorf("[buildSentenceLifecycleMsg]%v", err)
	}
	rtmSend.rtc.SendStreamMessage(msgBytes)
	logger.InfoContext(ctx, "[rtm] 发送session消息到rtm成功", zap.Int64("sgid", sgid), zap.Any("flag", flag))
	return nil
}

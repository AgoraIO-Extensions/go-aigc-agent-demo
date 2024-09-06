package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/pkg/logger"
)

func (tts *TTS) synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	//logger.Inst().Info("[tts] tts sentence Started!")
}

func (tts *TTS) synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	//logger.Inst().Info(fmt.Sprintf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData)))
}

func (tts *TTS) synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	// event.Result.AudioData 这里面存储的是tts识别的语音全部数据，且包含了wav的头部字节（44+2）
	//logger.Inst().Info(fmt.Sprintf("[tts] 音频总长度：%d，音频时间长度：%.2fs\n", len(event.Result.AudioData), event.Result.AudioDuration.Seconds()))
	//fmt.Printf("event.Result.AudioData[0:46]:%v\n", event.Result.AudioData[:46])
	//fmt.Printf("event.Result.AudioData[46:96]:%v\n", event.Result.AudioData[46:94])
}

func (tts *TTS) cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	logger.Error(fmt.Sprintf("[tts] 语音合成触发cancel. Reason:%v", event.Result.Reason))
}

package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/pkg/logger"
)

func (tts *TTS) synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
}

func (tts *TTS) synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
}

func (tts *TTS) synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
}

func (tts *TTS) cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	logger.ErrorContext(tts.ctx, fmt.Sprintf("[tts cancelledHandler] Reason:%v", event.Result.Reason))
}

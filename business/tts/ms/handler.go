package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/pkg/logger"
)

func (s *speechSynthesizer) synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
}

func (s *speechSynthesizer) synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
}

func (s *speechSynthesizer) synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	logger.DebugContext(s.ctx, "[tts][synthesizedHandler]")
	s.ctx = nil
	pool.put(s)
}

func (s *speechSynthesizer) cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	logger.ErrorContext(s.ctx, fmt.Sprintf("[tts][cancelledHandler] Reason:%v", event.Result.Reason))
}

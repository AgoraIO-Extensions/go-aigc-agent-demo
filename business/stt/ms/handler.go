package ms

import (
	"fmt"
	ms_common "github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"strings"
)

func (c *client) sessionStartedHandler(event speech.SessionEventArgs) {
	defer event.Close()
}

func (c *client) sessionStoppedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	logger.InfoContext(c.ctx, fmt.Sprintf("[stt sessionStoppedHandler] text list：%+v", c.results))
	c.result <- &common.Result{Text: strings.Join(c.results, ""), Complete: true}
	c.stop <- struct{}{}
}

func (c *client) recognizingHandler(event speech.SpeechRecognitionEventArgs) {
	text := event.Result.Text
	combineText := strings.Join(c.results, "")
	c.result <- &common.Result{Text: combineText + text}
	defer event.Close()
}

func (c *client) recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	/* STT may internally segment the audio, so the recognition results need to be saved here and then merged together in the end */
	c.results = append(c.results, event.Result.Text)
	c.result <- &common.Result{Text: strings.Join(c.results, "")}
}

func (c *client) cancelledHandler(event speech.SpeechRecognitionCanceledEventArgs) {
	defer event.Close()
	switch event.Reason {
	case ms_common.Error:
		logger.ErrorContext(c.ctx, fmt.Sprintf("[stt cancelledHandler] error event. ErrorDetails:%s, reason:%s", event.ErrorDetails, event.Reason))
		c.result <- &common.Result{Fail: true}
	case ms_common.CancelledByUser:
		logger.InfoContext(c.ctx, fmt.Sprintf("[stt cancelledHandler] ErrorDetails:%s, reason:%s", event.ErrorDetails, event.Reason))
		c.result <- &common.Result{Fail: true}
	case ms_common.EndOfStream:
	}
}

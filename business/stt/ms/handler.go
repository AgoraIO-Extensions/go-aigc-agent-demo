package ms

import (
	"fmt"
	ms_common "github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"strings"
)

func (c *client) sessionStartedHandler(event speech.SessionEventArgs) {
	defer event.Close()
}

func (c *client) sessionStoppedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	logger.Info(fmt.Sprintf("[stt stop回调] 文本列表：%+v", c.results), slog.Int64("sid", c.sid))
	c.result <- &common.Result{Text: strings.Join(c.results, ""), Complete: true}
	c.stop <- struct{}{}
}

// recognizingHandler 对一些过程中的临时句子进行处理
func (c *client) recognizingHandler(event speech.SpeechRecognitionEventArgs) {
	text := event.Result.Text
	combineText := strings.Join(c.results, "")
	c.result <- &common.Result{Text: combineText + text}
	defer event.Close()
}

// recognizedHandler 对最终的句子进行处理
func (c *client) recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	/* 尽管上游filter对音频进行了断句，但stt依然可能对输入的音频再次断句，因此这里需要将stt返回的 1~n 个文本临时保存起来，最后会合并成一个句子 */
	c.results = append(c.results, event.Result.Text)
	c.result <- &common.Result{Text: strings.Join(c.results, "")}
}

func (c *client) cancelledHandler(event speech.SpeechRecognitionCanceledEventArgs) {
	defer event.Close()
	switch event.Reason {
	case ms_common.Error:
		logger.Error(fmt.Sprintf("[stt cancel回调] 触发cancellation事件错误. ErrorDetails:%s, reason:%s", event.ErrorDetails, event.Reason), slog.Int64("sid", c.sid))
		c.result <- &common.Result{Fail: true}
	case ms_common.CancelledByUser:
		logger.Info(fmt.Sprintf("[stt cancel回调] 触发cancellation事件. ErrorDetails:%s, reason:%s", event.ErrorDetails, event.Reason), slog.Int64("sid", c.sid))
		c.result <- &common.Result{Fail: true}
	case ms_common.EndOfStream:
	}
}

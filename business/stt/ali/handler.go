package ali

import (
	"encoding/json"
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type response struct {
	Header  header  `json:"header"`
	Payload payload `json:"payload"`
}

type header struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
	MessageID  string `json:"message_id"`
	TaskID     string `json:"task_id"`
	StatusText string `json:"status_text"`
}

type stashResult struct {
	SentenceID  int           `json:"sentenceId"`
	BeginTime   int           `json:"beginTime"`
	Text        string        `json:"text"`
	FixedText   string        `json:"fixedText"`
	UnfixedText string        `json:"unfixedText"`
	CurrentTime int           `json:"currentTime"`
	Words       []interface{} `json:"words"`
}

type payload struct {
	Index      int     `json:"index"`
	Time       int     `json:"time"`
	Result     string  `json:"result"`
	BeginTime  int     `json:"begin_time"`
	Confidence float64 `json:"confidence"`
	//Words      []interface{} `json:"words"`
	//Status     int           `json:"status"`
	//Gender     string        `json:"gender"`
	//FixedResult string       `json:"fixed_result"`
	//UnfixedResult string     `json:"unfixed_result"`
	//StashResult stashResult  `json:"stash_result"`
	//AudioExtraInfo string    `json:"audio_extra_info"`
	//SentenceID string        `json:"sentence_id"`
	//GenderScore float64      `json:"gender_score"`
}

func (c *conn) onTaskFailed(jsonStr string, _ interface{}) {
	if c.ctx != nil {
		logger.ErrorContext(c.ctx, fmt.Sprintf("[onTaskFailed]:%s", jsonStr))
	}
	c.result <- &common.Result{Fail: true}

}

// onStarted This function is triggered by the Start function
func (c *conn) onStarted(jsonStr string, _ interface{}) {
	//logger.Inst().Info(fmt.Sprintf("onStarted:%s", text), slog.Int64("sid", c.sid))
}

func (c *conn) onSentenceBegin(jsonStr string, _ interface{}) {
	//logger.Inst().Info(fmt.Sprintf("onSentenceBegin:%s", text), slog.Int64("sid", c.sid))
}

// onResultChanged Return the Intermediate recognition result
func (c *conn) onResultChanged(jsonStr string, _ interface{}) {
	var err error
	defer func() {
		if err != nil {
			c.result <- &common.Result{Fail: true}
		}
	}()

	resp := response{}
	if err = json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		logger.ErrorContext(c.ctx, "[onResultChanged] Failed to unmarshal the text JSON returned by ali-stt", slog.Any("err", err))
		return
	}
	text := resp.Payload.Result
	logger.InfoContext(c.ctx, "[onResultChanged] Intermediate values", slog.String("text", text))
	c.result <- &common.Result{Text: text}
}

/*
onSentenceEnd: Return the final recognition result
Function trigger conditions：
 1. Either when the audio silence duration reaches the MaxSentenceSilence value, the server triggers it.
 2. Or when the client actively executes the stop function to trigger it.

Note: If no silence packets are sent after completing a sentence and the stop function is not actively called, this function will never be triggered.
*/
func (c *conn) onSentenceEnd(jsonStr string, _ interface{}) {
	resp := response{}
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		logger.ErrorContext(c.ctx, "[onSentenceEnd]Failed to unmarshal the text JSON returned by STT", slog.Any("err", err))
		c.result <- &common.Result{Fail: true}
		return
	}
	c.returnedAns = true
	text := resp.Payload.Result
	logger.InfoContext(c.ctx, "[onSentenceEnd] recognition result", slog.String("text", text))
	c.result <- &common.Result{Text: text, Complete: true}
	close(c.result)
}

func (c *conn) onCompleted(jsonStr string, _ interface{}) {
	logger.InfoContext(c.ctx, fmt.Sprintf("[stt] onCompleted:%s", jsonStr))
	if !c.returnedAns {
		c.result <- &common.Result{Text: "", Complete: true}
		close(c.result)
	}
}

func (c *conn) onClose(_ interface{}) {
	if c.ctx != nil {
		logger.InfoContext(c.ctx, fmt.Sprintf("onClose"))
	}
}

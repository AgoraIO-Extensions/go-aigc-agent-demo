package ali

import (
	"encoding/json"
	"fmt"
	"go-aigc-agent-demo/business/sentencelifecycle"
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
	if c.sid != 0 {
		logger.Info(fmt.Sprintf("onTaskFailed:%s", jsonStr), slog.Int64("sid", c.sid))
		c.result <- &common.Result{Fail: true}
	}
}

// 该函数由 Start 函数触发
func (c *conn) onStarted(jsonStr string, _ interface{}) {
	//logger.Inst().Info(fmt.Sprintf("onStarted:%s", text), slog.Int64("sid", c.sid))
}

func (c *conn) onSentenceBegin(jsonStr string, _ interface{}) {
	//logger.Inst().Info(fmt.Sprintf("onSentenceBegin:%s", text), slog.Int64("sid", c.sid))
}

// 返回过程中的识别结果
func (c *conn) onResultChanged(jsonStr string, _ interface{}) {
	var err error
	defer func() {
		if err != nil {
			c.result <- &common.Result{Fail: true}
		}
	}()

	resp := response{}
	if err = json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		logger.Error("ali stt返回文本json unmarshal失败", slog.Any("err", err), slog.Int64("sid", c.sid))
		return
	}
	text := resp.Payload.Result
	logger.Debug("[stt] onResultChanged 识别到文本中间值", slog.Int64("sid", c.sid), slog.String("text", text))
	c.result <- &common.Result{Text: text}
}

/*
onSentenceEnd: 返回最终的识别结果
函数触发条件：
 1. 要么音频静默时长达到 MaxSentenceSilence 值，服务端触发。
 2. 要么客户端主动执行stop函数触发。

注意：如果发完一句话之后不发静音包也不主动调用Stop函数，那么这个函数永远不会被触发
*/
func (c *conn) onSentenceEnd(jsonStr string, _ interface{}) {
	resp := response{}
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		logger.Error("ali stt返回文本json unmarshal失败", slog.Any("err", err), slog.Int64("sid", c.sid))
		c.result <- &common.Result{Fail: true}
		return
	}
	c.returnedAns = true
	text := resp.Payload.Result
	logger.Info("[stt] onSentenceEnd 识别到文本", sentencelifecycle.Tag(c.sid), slog.String("text", text))
	c.result <- &common.Result{Text: text, Complete: true}
	close(c.result)
}

func (c *conn) onCompleted(jsonStr string, _ interface{}) {
	logger.Info(fmt.Sprintf("[stt] onCompleted:%s", jsonStr), slog.Int64("sid", c.sid))
	if !c.returnedAns {
		c.result <- &common.Result{Text: "", Complete: true}
		close(c.result)
	}
}

func (c *conn) onClose(_ interface{}) {
	if c.sid != 0 {
		logger.Info(fmt.Sprintf("onClose"), slog.Int64("sid", c.sid))
	}
}

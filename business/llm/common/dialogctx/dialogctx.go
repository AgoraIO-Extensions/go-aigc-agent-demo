package dialogctx

import (
	"sync"
)

type Role string

const (
	SYSTEM    = "system"
	USER      = "user"
	ASSISTANT = "assistant"
)

// qa Q&A
type qa struct {
	qid int64  // 问题id（id必须是递增的）
	q   string // 问题内容
	a   string // 回答内容
}

type Message struct {
	qid     int64  // 问题id或者是合并后的问题id（等于合并问题中最新的问题的id）； 注意：该字段小写非导出（不会被JSON序列化）
	Role    Role   `json:"role"`
	Content string `json:"content"` // 问题内容或者是合并后的问题内容
}

/* ------------------------------------------------------------------------------------------------------------------ */

type DialogCTX struct {
	WithHistory bool
	qaList      []*qa
	qaMap       *sync.Map // key: qid value: *qa
	maxQNum     int
	latestQID   int64 // 最新的 qid
}

func NewDialogCTX(nums int, WithHistory bool) *DialogCTX {
	if nums <= 0 {
		nums = 1
	}
	return &DialogCTX{
		WithHistory: WithHistory,
		qaList:      make([]*qa, 0),
		qaMap:       new(sync.Map),
		maxQNum:     nums,
	}
}

// AddQuestion 构建上下文信息链。（此函数并发不安全）
func (dCtx *DialogCTX) AddQuestion(question string, sgid int64) []Message {
	qid := sgid
	dCtx.latestQID = qid

	if !dCtx.WithHistory {
		return []Message{{
			qid:     qid,
			Role:    USER,
			Content: question,
		}}
	}

	// 更新 qaList 和 qaMap
	newUint := &qa{qid: qid, q: question}
	n := len(dCtx.qaList)
	if _, ok := dCtx.qaMap.Load(qid); ok {
		dCtx.qaList[n-1] = newUint
	} else {
		dCtx.qaList = append(dCtx.qaList, newUint)
	}
	if len(dCtx.qaList) > dCtx.maxQNum {
		deleteQid := dCtx.qaList[0].qid
		dCtx.qaList = dCtx.qaList[1:]
		dCtx.qaMap.Delete(deleteQid)
	}
	dCtx.qaMap.Store(qid, newUint)

	// 转换为 []Message
	msgs := make([]Message, 0, len(dCtx.qaList))
	for _, qaUint := range dCtx.qaList {
		msgs = append(msgs,
			Message{
				qid:     qaUint.qid,
				Role:    USER,
				Content: qaUint.q,
			}, Message{
				qid:     qaUint.qid,
				Role:    ASSISTANT,
				Content: qaUint.a,
			})
	}

	return msgs
}

// AddAnswer 给最新的问题 添加/追加补充 答案
// 返回true表示当前goroutine无需继续提供流式输出的答案了
func (dCtx *DialogCTX) AddAnswer(ansSegment string, sgid int64) bool {
	questionID := sgid
	if !dCtx.WithHistory {
		return false
	}
	unitAny, ok := dCtx.qaMap.Load(questionID)
	if !ok {
		return true
	}
	dialog := unitAny.(*qa)
	dialog.a = dialog.a + ansSegment
	return false
}

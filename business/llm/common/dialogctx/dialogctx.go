package dialogctx

import (
	"fmt"
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
	qid int64  // Question ID (ID must be incremental)
	q   string // Question content
	a   string // Answer content
}

type Message struct {
	qid     int64  // Question ID or Question ID after merge
	Role    Role   `json:"role"`
	Content string `json:"content"` // Question content or Question content after merge
}

/* ------------------------------------------------------------------------------------------------------------------ */

type DialogCTX struct {
	WithHistory bool
	qaList      []*qa
	qaMap       *sync.Map // key: qid value: *qa
	maxQNum     int
	latestQID   int64
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

// AddQuestion Build a context information chain. (concurrency unsafe)"
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

	// update qaList & qaMap
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

	// format []*qa to []Message
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

// StreamAddAnswer added/appended answer streamingly
func (dCtx *DialogCTX) StreamAddAnswer(ansSegment string, sgid int64) error {
	questionID := sgid
	if !dCtx.WithHistory {
		return nil
	}
	unitAny, ok := dCtx.qaMap.Load(questionID)
	if !ok {
		return fmt.Errorf("the QA pair corresponding to sgid(%d) could not be found in dCtx.qaMap", sgid)
	}
	dialog := unitAny.(*qa)
	dialog.a = dialog.a + ansSegment
	return nil
}

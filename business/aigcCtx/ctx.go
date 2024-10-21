package aigcCtx

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/aigcCtx/sentence"
	"go-aigc-agent-demo/pkg/logger"
	"sync"
	"time"
)

var head *AIGCContext
var tail *AIGCContext
var locker = new(sync.Mutex)

func init() {
	head = new(AIGCContext)
	tail = new(AIGCContext)
	head.next = tail
	tail.prev = head
}

type AIGCContext struct {
	context.Context
	prev     *AIGCContext
	next     *AIGCContext
	MetaData *sentence.MetaData
	cancel   context.CancelFunc
	canceled bool
}

func NewContext(pCtx context.Context, metaData *sentence.MetaData) *AIGCContext {
	locker.Lock()
	defer locker.Unlock()
	ctx, cancel := context.WithCancel(pCtx)

	node := &AIGCContext{
		prev:     tail.prev,
		next:     tail,
		MetaData: metaData,
		Context:  ctx,
		cancel:   cancel,
	}
	tail.prev.next = node
	tail.prev = node
	return node
}

func (ctx *AIGCContext) ReleaseCtxNode() {
	locker.Lock()
	defer locker.Unlock()
	if ctx.canceled {
		return
	}
	ctx.cancel()
	ctx.canceled = true
	logger.InfoContext(ctx.Context, fmt.Sprintf("[interrupt] sid:%d will be released", ctx.MetaData.Sid))
	ctx.prev.next = ctx.next
	ctx.next.prev = ctx.prev
}

// WaitNodesCancel 等待此刻存在的后续节点释放
func (ctx *AIGCContext) WaitNodesCancel() <-chan struct{} {
	done := make(chan struct{}, 1)
	maxSid := tail.prev.MetaData.Sid

	go func() {
		for {
			if ctx.next == tail || ctx.next.MetaData.Sid > maxSid {
				done <- struct{}{}
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()

	return done
}

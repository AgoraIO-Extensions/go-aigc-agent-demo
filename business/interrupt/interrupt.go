package interrupt

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
	"sync"
)

var head *CtxNode
var tail *CtxNode
var locker = new(sync.Mutex)

func init() {
	head = new(CtxNode)
	tail = new(CtxNode)
	head.Sid = -1
	head.next = tail
	tail.Sid = -2
	tail.prev = head
}

type CtxNode struct {
	prev     *CtxNode
	next     *CtxNode
	Sid      int64
	Ctx      context.Context
	cancel   context.CancelFunc
	canceled bool
}

func NewCtxNode(pCtx context.Context, sid int64) *CtxNode {
	locker.Lock()
	defer locker.Unlock()
	ctx, cancel := context.WithCancel(pCtx)

	node := &CtxNode{
		prev:   tail.prev,
		next:   tail,
		Sid:    sid,
		Ctx:    ctx,
		cancel: cancel,
	}
	tail.prev.next = node
	tail.prev = node
	return node
}

func (node *CtxNode) ReleaseCtxNode() {
	locker.Lock()
	defer locker.Unlock()
	if node.canceled {
		return
	}
	node.cancel()
	node.canceled = true
	logger.InfoContext(node.Ctx, fmt.Sprintf("[interrupt] sid:%d will be released", node.Sid))
	node.prev.next = node.next
	node.next.prev = node.prev
}

// Interrupt node interrupt node.prev
func (node *CtxNode) Interrupt() {
	locker.Lock()
	defer locker.Unlock()

	if node.canceled {
		return
	}
	nd := node.prev
	for nd != head {
		logger.InfoContext(node.Ctx, fmt.Sprintf("[interrupt] sid:%d will interrupt sid:%d", node.Sid, nd.Sid))
		nd.cancel()
		nd.canceled = true
		nd = nd.prev
	}

	head.next = node
	node.prev = head
}

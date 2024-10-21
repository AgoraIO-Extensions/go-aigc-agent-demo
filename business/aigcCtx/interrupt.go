package aigcCtx

import (
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
)

// Interrupt node interrupt node.prev
func (ctx *AIGCContext) Interrupt() {
	locker.Lock()
	defer locker.Unlock()

	if ctx.canceled {
		return
	}
	nd := ctx.prev
	for nd != head {
		logger.InfoContext(ctx, fmt.Sprintf("[interrupt] sid:%d will interrupt sid:%d", ctx.MetaData.Sid, nd.MetaData.Sid))
		nd.cancel()
		nd.canceled = true
		nd = nd.prev
	}

	head.next = ctx
	ctx.prev = head
}

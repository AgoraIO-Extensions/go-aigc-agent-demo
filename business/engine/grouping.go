package engine

import (
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/config"
	"time"
)

// Grouping set sgid for sentence with sid {sid}
func Grouping(oldSgid, sid int64, prevSentenceEndTime time.Time) int64 {
	cfg := config.Inst()
	sgid := oldSgid

	if cfg.GroupingStrategy == config.DependOnRTCSend {
		if sentencelifecycle.IfSidIntoRTC(sid - 1) {
			sgid = sid
			sentencelifecycle.DeleteSidIntoRtc(sid - 1)
		}
	}

	if cfg.GroupingStrategy == config.DependOnTime {
		if time.Since(prevSentenceEndTime) > time.Millisecond*500 {
			sgid = sid
		}
	}

	return sgid
}

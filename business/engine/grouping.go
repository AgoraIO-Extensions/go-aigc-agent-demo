package engine

import (
	"go-aigc-agent-demo/business/aigcCtx/sentence"
	"go-aigc-agent-demo/config"
	"time"
)

// Grouping set sgid for sentence with sid {sid}
func Grouping(prevSMetaData *sentence.MetaData, sid int64) int64 {
	cfg := config.Inst()
	sgid := prevSMetaData.Sgid
	if sgid == 0 {
		return sid
	}

	if cfg.Grouping.Strategy == config.DependOnRTCSend {
		if prevSMetaData.StageSendToRTC {
			sgid = sid
		}
	}

	if cfg.Grouping.Strategy == config.DependOnTime {
		if time.Since(prevSMetaData.FilterAudioTailRcvTime).Milliseconds() > config.Inst().Grouping.TimeThreshold {
			sgid = sid
		}
	}

	return sgid
}

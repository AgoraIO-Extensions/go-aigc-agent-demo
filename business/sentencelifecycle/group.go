package sentencelifecycle

import (
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

func GroupInst() *Groups {
	return groups
}

var groups = &Groups{
	sidTosgid:       new(sync.Map),
	sgidToStartTime: make(map[int64]*time.Time),
}

type Groups struct {
	sidTosgid       *sync.Map
	sgidToStartTime map[int64]*time.Time
}

/*  ------------------------------------------------- sid ——> sgid ------------------------------------------------- */

func (g *Groups) SetSidToSgid(sid int64, sgid int64) {
	g.sidTosgid.Store(sid, sgid)
}

func (g *Groups) GetSgidBySid(sid int64) int64 {
	v, ok := g.sidTosgid.Load(sid)
	if !ok {
		logger.Inst().Error("此处为理论不可达代码，出现了代表代码存在bug", zap.Int64("sid", sid))
		return 0
	}
	return v.(int64)
}

func (g *Groups) DeleteSidToSgid(sid int64) {
	g.sidTosgid.Delete(sid)
}

/* ---------------------------------------------- sgid ——> sentencelifecycle group内输入音频的结束时间 ---------------------------------------------- */

// StoreInAudioEndTimeInOneSentenceGroup 存储一个sentence group中，来自rtc输入音频的最后一个chunk的接收时间
func (g *Groups) StoreInAudioEndTimeInOneSentenceGroup(sgid int64, startTime time.Time) {
	g.sgidToStartTime[sgid] = &startTime
}

func (g *Groups) GetInAudioEndTimeInOneSentenceGroup(sgid int64) *time.Time {
	return g.sgidToStartTime[sgid]
}

func (g *Groups) DeleteInAudioEndTimeInOneSentenceGroup(sgid int64) {
	delete(g.sgidToStartTime, sgid)
}

/* ----------------------------------------------------- log tag ---------------------------------------------------- */

func Tag(sid int64, sgid ...int64) zap.Field {
	if len(sgid) > 0 {
		return zap.Int64s("sid,sgid", []int64{sid, sgid[0]})
	}
	v, ok := groups.sidTosgid.Load(sid)
	if !ok {
		return zap.Int64("sid", sid)
	}
	return zap.Int64s("sid,sgid", []int64{sid, v.(int64)})
}

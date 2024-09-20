package sentencelifecycle

import (
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

func GroupInst() *Groups {
	return groups
}

var groups = &Groups{
	sidTosgid:      new(sync.Map),
	sidToStartTime: make(map[int64]*time.Time),
}

type Groups struct {
	sidTosgid      *sync.Map
	sidToStartTime map[int64]*time.Time
}

/*  ------------------------------------------------- sid ——> sgid ------------------------------------------------- */

func (g *Groups) SetSidToSgid(sid int64, sgid int64) {
	g.sidTosgid.Store(sid, sgid)
}

func (g *Groups) GetSgidBySid(sid int64) int64 {
	v, ok := g.sidTosgid.Load(sid)
	if !ok {
		logger.Error("此处为理论不可达代码，出现了代表代码存在bug", slog.Int64("sid", sid))
		return 0
	}
	return v.(int64)
}

func (g *Groups) DeleteSidToSgid(sid int64) {
	g.sidTosgid.Delete(sid)
}

/* ---------------------------------------------- sid —> end time of a sentence input audio ---------------------------------------------- */

// StoreAudioEndTime Store the reception time of the last chunk of input audio from RTC within a sentence
func (g *Groups) StoreAudioEndTime(sid int64, startTime time.Time) {
	g.sidToStartTime[sid] = &startTime
}

func (g *Groups) GetAudioEndTime(sid int64) *time.Time {
	return g.sidToStartTime[sid]
}

func (g *Groups) DeleteAudioEndTime(sid int64) {
	delete(g.sidToStartTime, sid)
}

/* ----------------------------------------------------- log tag ---------------------------------------------------- */

func Tag(sid int64, sgid ...int64) slog.Attr {
	if len(sgid) > 0 {
		return slog.Any("sid,sgid", []int64{sid, sgid[0]})
	}
	v, ok := groups.sidTosgid.Load(sid)
	if !ok {
		return slog.Int64("sid", sid)
	}
	return slog.Any("sid,sgid", []int64{sid, v.(int64)})
}

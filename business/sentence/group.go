package sentence

import (
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

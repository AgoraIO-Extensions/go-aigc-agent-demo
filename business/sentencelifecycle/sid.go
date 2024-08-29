package sentencelifecycle

import (
	"time"
)

/* --------------------------------------------------- 初始化sid --------------------------------------------------- */

var FirstSid = time.Now().Unix()

/* ----------------------------------------------- 记录/判断sid对应的音频是否发到过rtc ----------------------------------------------- */

var sidToRTC = make(map[int64]bool)

func SetSidIntoRTC(sid int64) {
	sidToRTC[sid] = true
}

func IfSidIntoRTC(sid int64) bool {
	return sidToRTC[sid]
}

func DeleteSidIntoRtc(sid int64) {
	delete(sidToRTC, sid)
}

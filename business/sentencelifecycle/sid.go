package sentencelifecycle

import (
	"time"
)

/* --------------------------------------------------- Initialize sid --------------------------------------------------- */

var FirstSid = time.Now().Unix()

/* ----------------------------------------------- Record/judge whether the audio corresponding to the sid has been sent to RTC ----------------------------------------------- */

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

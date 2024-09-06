package common

import (
	"time"
)

type Segment struct {
	AudioChan chan []byte // 每个chunk 320 bytes
	Sid       int64
	ID        int
	Text      string
	SendTime  time.Time // 发送到tts的时间
}

type Sentence struct {
	ID        int64
	SegChan   chan *Segment
	AudioChan chan []byte // 每个chunk 320 bytes
}

func (s *Sentence) mergeSegments() {
	go func() {
		defer close(s.AudioChan)
		for {
			seg, ok := <-s.SegChan
			if !ok {
				return
			}
			for {
				aud, ok := <-seg.AudioChan
				if !ok {
					break
				}
				s.AudioChan <- aud
			}
		}
	}()
}

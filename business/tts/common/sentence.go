package common

import (
	"time"
)

type Segment struct {
	AudioChan chan []byte // 320 bytes per chunk
	//Sid       int64
	ID       int
	Text     string
	SendTime time.Time // time sent to TTS
}

type Sentence struct {
	//ID        int64
	SegChan   chan *Segment
	AudioChan chan []byte // 320 bytes per chunk
}

func (s *Sentence) mergeSegments() {
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
}

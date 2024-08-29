package sentence

import "time"

type Sentence struct {
	SID           int64
	SegChan       <-chan string
	BeginTime     time.Time
	SegInChanTime time.Time
	IsExp         bool
	Failed        bool
}

func (s *Sentence) SetSegChan(segChan <-chan string) {
	s.SegChan = segChan
}

func (s *Sentence) WaitSegChan(dur time.Duration) bool {
	for {
		select {
		case <-time.After(dur):
			return false
		default:
			if s.Failed {
				return false
			}
			if s.SegChan != nil {
				return true
			}
			time.Sleep(time.Millisecond * 2)
		}
	}
}

type SentenceChan chan *Sentence

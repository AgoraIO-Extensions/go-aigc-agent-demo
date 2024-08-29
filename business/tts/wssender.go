package tts

import (
	"fmt"
)

type webSocketTTS interface {
	Send(sid int64, text string) (<-chan byte, error) // 返回当前session的音频流，只会在第一次发送时返回
}

type webSocketSender struct {
	client webSocketTTS
}

// send text为空表示发送结束
func (w *webSocketSender) send(ss *Sentence, text string) error {
	audioChan, err := w.client.Send(ss.ID, text)
	if err != nil {
		return fmt.Errorf("[w.client.SendPcm]%v", err)
	}

	if audioChan != nil {
		go func(audioChan <-chan byte) {
			defer close(ss.AudioChan)
			var ok bool
			for {
				chunk := make([]byte, 320)
				for i := range chunk {
					if chunk[i], ok = <-audioChan; !ok {
						return
					}
				}
				ss.AudioChan <- chunk
			}
		}(audioChan)
	}
	return nil
}

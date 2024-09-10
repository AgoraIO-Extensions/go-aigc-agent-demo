package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

func newPushAudioInputStream() (*audio.PushAudioInputStream, error) {
	stream, err := audio.CreatePushAudioInputStream()
	if err != nil {
		return nil, fmt.Errorf("[audio.CreatePushAudioInputStreamFromFormat]%v", err)
	}
	return stream, nil
}

// pumpChunkIntoStream write audio to stream
func (c *client) pumpChunkIntoStream(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	if err := c.sttInputStream.Write(buf[:]); err != nil {
		return fmt.Errorf("[stream.Write]%v", err)
	}
	return nil
}

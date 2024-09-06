package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

// newPushAudioInputStream 构建一个「stt输入流」
func newPushAudioInputStream() (*audio.PushAudioInputStream, error) {
	stream, err := audio.CreatePushAudioInputStream()
	if err != nil {
		return nil, fmt.Errorf("[audio.CreatePushAudioInputStreamFromFormat]%v", err)
	}
	return stream, nil
}

// pumpChunkIntoStream 将音频写到「stt输入流」中
func (c *client) pumpChunkIntoStream(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	if err := c.sttInputStream.Write(buf[:]); err != nil {
		return fmt.Errorf("[stream.Write]%v", err)
	}
	return nil
}

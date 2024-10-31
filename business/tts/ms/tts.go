package ms

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	ttscommon "go-aigc-agent-demo/business/tts/common"
	"io"
)

/* -------------------------------------------------- config ---------------------------------------------------------- */

type Config struct {
	setLog                      bool
	languageCheckMode           LanguageCheckMode
	specifyLanguage             string // Output audio language. Reference link: https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/language-support?tabs=tts
	outputVoice                 string // Output audio language and accent. Reference link: Same as above
	speechKey                   string
	speechRegion                string
	SpeechSynthesisOutputFormat common.SpeechSynthesisOutputFormat // Output audio format.
}

type LanguageCheckMode int

const (
	AutoCheck LanguageCheckMode = 0 // Automatic language detection mode
	Specify   LanguageCheckMode = 1 // Specify language mode
)

func Init(setLog bool, speechKey, speechRegion string, languageCheckMode int, specifyLanguage, outputVoice string, outputFormat common.SpeechSynthesisOutputFormat) error {
	cfg := &Config{
		setLog:                      setLog,
		languageCheckMode:           LanguageCheckMode(languageCheckMode),
		specifyLanguage:             specifyLanguage,
		outputVoice:                 outputVoice,
		speechKey:                   speechKey,
		speechRegion:                speechRegion,
		SpeechSynthesisOutputFormat: outputFormat,
	}

	speechConfig, err := speech.NewSpeechConfigFromSubscription(cfg.speechKey, cfg.speechRegion)
	if err != nil {
		return fmt.Errorf("[NewSpeechConfigFromSubscription]%w", err)
	}
	if err = speechConfig.SetSpeechSynthesisOutputFormat(cfg.SpeechSynthesisOutputFormat); err != nil {
		return fmt.Errorf("[SetSpeechSynthesisOutputFormat]%w", err)
	}
	if cfg.setLog {
		if err = speechConfig.SetProperty(common.SpeechLogFilename, "tts.log"); err != nil {
			return fmt.Errorf("[setLog]%v", err)
		}
	}

	initSynthesizerPool(cfg, speechConfig, 5)
	return nil
}

/* --------------------------------------------------- TTS --------------------------------------------------------- */

type TTS struct {
	*ttscommon.Sender
	ctx context.Context
}

func NewTTS(ctx context.Context, con int) *TTS {
	tts := &TTS{
		ctx: ctx,
	}
	tts.Sender = ttscommon.NewHttpSender(ctx, tts.streamAsk, con)
	return tts
}

// streamAsk send text to tts
func (tts *TTS) streamAsk(ctx context.Context, segID int, seg string) (io.ReadCloser, error) {
	syn, err := pool.get(ctx, segID)
	if err != nil {
		return nil, fmt.Errorf("[pool.get]%v", err)
	}
	task := syn.msSpeechSynthesizer.StartSpeakingTextAsync(seg)

	var outcome speech.SpeechSynthesisOutcome

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("interrupted while waiting for TTS to return, %w", context.Canceled)
	case outcome = <-task: // Return basic information of the server's speech recognition result (at this point, it does not include audio data)
	}

	defer outcome.Close()
	if outcome.Error != nil {
		return nil, fmt.Errorf("error occurred on the TTS server during the text-to-speech synthesis process：%v", outcome.Error)
	}
	stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
	if err != nil {
		return nil, fmt.Errorf("[NewAudioDataStreamFromSpeechSynthesisResult]%v", err)
	}

	return &streamReaderCloser{stream: stream}, nil
}

/* ---------------------------------------------------- streamReaderCloser -------------------------------------------------------- */

type streamReaderCloser struct {
	stream *speech.AudioDataStream
}

func (s *streamReaderCloser) Read(chunk []byte) (int, error) {
	return s.stream.Read(chunk)
}

func (s *streamReaderCloser) Close() error {
	s.stream.Close()
	return nil
}

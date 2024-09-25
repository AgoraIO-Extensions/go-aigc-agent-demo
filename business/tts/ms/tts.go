package ms

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	ttscommon "go-aigc-agent-demo/business/tts/common"
	"io"
)

func newSpeechSynthesizer(cfg *Config, speechConfig *speech.SpeechConfig) (*speech.SpeechSynthesizer, error) {
	if cfg.languageCheckMode == AutoCheck {
		langConfig, err := speech.NewAutoDetectSourceLanguageConfigFromOpenRange()
		if err != nil {
			return nil, fmt.Errorf("[NewAutoDetectSourceLanguageConfigFromLanguages]%v", err)
		}
		syn, err := speech.NewSpeechSynthesizerFomAutoDetectSourceLangConfig(speechConfig, langConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("[NewSpeechSynthesizerFomAutoDetectSourceLangConfig]%v", err)
		}
		return syn, nil
	}

	// If both SpeechSynthesisVoiceName and SpeechSynthesisLanguage are not set, the default voice for en-US will be used.
	// If only SpeechSynthesisLanguage is set, the default voice for the specified locale will be used.
	// If both SpeechSynthesisVoiceName and SpeechSynthesisLanguage are set, the SpeechSynthesisLanguage setting will be ignored. The system will use the voice specified by SpeechSynthesisVoiceName.
	if cfg.specifyLanguage != "" {
		if err := speechConfig.SetSpeechSynthesisLanguage(cfg.specifyLanguage); err != nil {
			return nil, fmt.Errorf("[SetSpeechSynthesisLanguage]%v", err)
		}
	}
	if cfg.outputVoice != "" {
		if err := speechConfig.SetSpeechSynthesisVoiceName(cfg.outputVoice); err != nil {
			return nil, fmt.Errorf("[SetSpeechSynthesisVoiceName]%v", err)
		}
	}

	syn, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("[NewSpeechSynthesizerFromConfig]%v", err)
	}
	return syn, nil
}

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

func NewTTSConfig(setLog bool, speechKey, speechRegion string, languageCheckMode int, specifyLanguage, outputVoice string, outputFormat common.SpeechSynthesisOutputFormat) *Config {
	return &Config{
		setLog:                      setLog,
		languageCheckMode:           LanguageCheckMode(languageCheckMode),
		specifyLanguage:             specifyLanguage,
		outputVoice:                 outputVoice,
		speechKey:                   speechKey,
		speechRegion:                speechRegion,
		SpeechSynthesisOutputFormat: outputFormat,
	}
}

/* --------------------------------------------------- TTS --------------------------------------------------------- */

type TTS struct {
	*ttscommon.HttpSender
	ctx               context.Context
	speechConfig      *speech.SpeechConfig
	speechSynthesizer *speech.SpeechSynthesizer
}

func NewTTS(ctx context.Context, con int, cfg *Config) (*TTS, error) {
	tts := &TTS{
		ctx: ctx,
	}
	tts.HttpSender = ttscommon.NewHttpSender(ctx, con, tts.streamAsk)
	var err error
	defer func() {
		if err != nil {
			tts.close()
		}
	}()

	tts.speechConfig, err = speech.NewSpeechConfigFromSubscription(cfg.speechKey, cfg.speechRegion)
	if err != nil {
		return nil, fmt.Errorf("[speech.NewSpeechConfigFromSubscription]%v", err)
	}

	if err = tts.speechConfig.SetSpeechSynthesisOutputFormat(cfg.SpeechSynthesisOutputFormat); err != nil {
		panic(err)
	}

	if cfg.setLog {
		if err = tts.speechConfig.SetProperty(common.SpeechLogFilename, "tts.log"); err != nil {
			return nil, fmt.Errorf("[setLog]%v", err)
		}
	}

	tts.speechSynthesizer, err = newSpeechSynthesizer(cfg, tts.speechConfig)
	if err != nil {
		return nil, fmt.Errorf("[newSpeechSynthesizer]%v", err)
	}

	tts.speechSynthesizer.SynthesisStarted(tts.synthesizeStartedHandler)
	tts.speechSynthesizer.Synthesizing(tts.synthesizingHandler)
	tts.speechSynthesizer.SynthesisCompleted(tts.synthesizedHandler)
	tts.speechSynthesizer.SynthesisCanceled(tts.cancelledHandler)

	return tts, nil
}

func (tts *TTS) close() {
	if tts.speechSynthesizer != nil {
		tts.speechSynthesizer.Close()
	}
	if tts.speechConfig != nil {
		tts.speechConfig.Close()
	}
}

// streamAsk send text to tts
func (tts *TTS) streamAsk(ctx context.Context, text string) (io.ReadCloser, error) {
	task := tts.speechSynthesizer.StartSpeakingTextAsync(text)

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

func (tts *TTS) Send(ctx context.Context, segmentID int, segmentContent string) {
	tts.HttpSender.Send(ctx, segmentID, segmentContent)
}

func (tts *TTS) GetResult() <-chan []byte {
	return tts.HttpSender.Result()
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

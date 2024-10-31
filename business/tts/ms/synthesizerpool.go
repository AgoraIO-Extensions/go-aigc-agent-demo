package ms

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"sync"
)

func newMSSpeechSynthesizer(cfg *Config, speechConfig *speech.SpeechConfig) (*speech.SpeechSynthesizer, error) {
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

type speechSynthesizer struct {
	ctx                 context.Context
	msSpeechSynthesizer *speech.SpeechSynthesizer
}

func newSpeechSynthesizer(cfg *Config, speechConfig *speech.SpeechConfig) (*speechSynthesizer, error) {
	msSyn, err := newMSSpeechSynthesizer(cfg, speechConfig)
	if err != nil {
		return nil, fmt.Errorf("[newMSSpeechSynthesizer]%w", err)
	}

	var syn = &speechSynthesizer{msSpeechSynthesizer: msSyn}
	msSyn.SynthesisStarted(syn.synthesizeStartedHandler)
	msSyn.Synthesizing(syn.synthesizingHandler)
	msSyn.SynthesisCompleted(syn.synthesizedHandler)
	msSyn.SynthesisCanceled(syn.cancelledHandler)

	return syn, nil
}

var pool *synthesizerPool

type synthesizerPool struct {
	cfg          *Config
	speechConfig *speech.SpeechConfig
	queue        chan *speechSynthesizer
	capacity     int
	locker       sync.Mutex
}

func initSynthesizerPool(cfg *Config, speechConfig *speech.SpeechConfig, capacity int) {
	pool = &synthesizerPool{
		cfg:          cfg,
		speechConfig: speechConfig,
		queue:        make(chan *speechSynthesizer, capacity),
		capacity:     capacity,
		locker:       sync.Mutex{},
	}
}

func (s *synthesizerPool) get(ctx context.Context, segID int) (*speechSynthesizer, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if len(s.queue) == 0 {
		logger.InfoContext(ctx, "[tts] will new a speechSynthesizer", slog.Int("segID", segID))
		return newSpeechSynthesizer(s.cfg, s.speechConfig)
	}

	syn := <-s.queue
	syn.ctx = ctx
	return syn, nil
}

func (s *synthesizerPool) put(synthesizer *speechSynthesizer) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if len(s.queue) == s.capacity {
		syn := <-s.queue
		syn.msSpeechSynthesizer.Close()
	}
	s.queue <- synthesizer
}

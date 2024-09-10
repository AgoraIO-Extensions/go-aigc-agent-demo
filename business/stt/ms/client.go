package ms

import (
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	ms_common "github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"time"
)

type client struct {
	audioConfig      *audio.AudioConfig
	speechConfig     *speech.SpeechConfig
	speechRecognizer *speech.SpeechRecognizer
	sttInputStream   *audio.PushAudioInputStream
	stop             chan struct{}
	sid              int64
	results          []string
	result           chan *common.Result
}

func newClient(sid int64, cfg *Config) (*client, error) {
	c := &client{
		stop:   make(chan struct{}, 1),
		sid:    sid,
		result: make(chan *common.Result, 100),
	}
	var err error
	defer func() {
		if err != nil {
			c.close()
		}
	}()

	/* Currently, only support (16 kHz, 16 bit, mono PCM)  */
	c.sttInputStream, err = newPushAudioInputStream() // 只支持 (16 kHz, 16 bit, mono PCM)
	if err != nil {
		return nil, fmt.Errorf("[newPushAudioInputStream]%v", err)
	}
	c.audioConfig, err = audio.NewAudioConfigFromStreamInput(c.sttInputStream)
	if err != nil {
		return nil, fmt.Errorf("[NewAudioConfigFromStreamInput]%v", err)
	}
	c.speechConfig, err = speech.NewSpeechConfigFromSubscription(cfg.speechKey, cfg.speechRegion)
	if err != nil {
		return nil, fmt.Errorf("[NewSpeechConfigFromSubscription]%v", err)
	}
	if cfg.setLog {
		if err = c.speechConfig.SetProperty(ms_common.SpeechLogFilename, "stt.log"); err != nil {
			return nil, fmt.Errorf("[setLog]%v", err)
		}
	}
	c.speechRecognizer, err = newSpeechRecognizer(cfg, c.speechConfig, c.audioConfig)
	if err != nil {
		return nil, fmt.Errorf("[newSpeechRecognizer]%v", err)
	}
	c.speechRecognizer.SessionStarted(c.sessionStartedHandler)
	c.speechRecognizer.SessionStopped(c.sessionStoppedHandler)
	c.speechRecognizer.Recognizing(c.recognizingHandler)
	c.speechRecognizer.Recognized(c.recognizedHandler)
	c.speechRecognizer.Canceled(c.cancelledHandler)
	c.speechRecognizer.StartContinuousRecognitionAsync()
	c.asyncCloseSTT()

	return c, nil
}

func (c *client) asyncCloseSTT() {
	go func() {
		select {
		case <-c.stop:
			c.close()
			logger.Info("[stt] release resource", sentencelifecycle.Tag(c.sid))
		case <-time.After(time.Second * 30):
			c.close()
			logger.Info("[stt] Exceeded 30 seconds without receiving the recognition end signal, will release resource immediately.", sentencelifecycle.Tag(c.sid))
		}
	}()
}

func newSpeechRecognizer(cfg *Config, speechConfig *speech.SpeechConfig, audioConfig *audio.AudioConfig) (*speech.SpeechRecognizer, error) {
	if cfg.languageCheckMode == AutoCheck {
		langConfig, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages(cfg.autoAudioCheckLanguage)
		if err != nil {
			return nil, fmt.Errorf("[NewAutoDetectSourceLanguageConfigFromLanguages]%v", err)
		}
		rec, err := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(speechConfig, langConfig, audioConfig)
		if err != nil {
			return nil, fmt.Errorf("[NewSpeechRecognizerFomAutoDetectSourceLangConfig]%v", err)
		}
		return rec, nil
	}

	if err := speechConfig.SetSpeechRecognitionLanguage(cfg.specifyLanguage); err != nil {
		return nil, fmt.Errorf("[SetSpeechRecognitionLanguage]%v", err)
	}
	rec, err := speech.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return nil, fmt.Errorf("[NewSpeechRecognizerFromConfig]%v", err)
	}
	return rec, nil
}

func (c *client) close() {
	if c.sttInputStream != nil {
		c.sttInputStream.CloseStream()
	}
	if c.speechRecognizer != nil {
		c.speechRecognizer.StopContinuousRecognitionAsync()
	}
	if c.speechRecognizer != nil {
		c.speechRecognizer.Close()
	}
	if c.speechConfig != nil {
		c.speechConfig.Close()
	}
	if c.audioConfig != nil {
		c.audioConfig.Close()
	}
	if c.sttInputStream != nil {
		c.sttInputStream.Close()
	}
}

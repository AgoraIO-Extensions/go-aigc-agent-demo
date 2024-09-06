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

	// 如果没有设置 SpeechSynthesisVoiceName 或 SpeechSynthesisLanguage，则会讲 en-US 的默认语音。
	// 如果仅设置了 SpeechSynthesisLanguage，则会讲指定区域设置的默认语音。
	// 如果同时设置了 SpeechSynthesisVoiceName 和 SpeechSynthesisLanguage，则会忽略 SpeechSynthesisLanguage 设置。 系统会讲你使用 SpeechSynthesisVoiceName 指定的语音。
	// 如果使用语音合成标记语言 (SSML) 设置了 voice 元素，则会忽略 SpeechSynthesisVoiceName 和 SpeechSynthesisLanguage 设置。
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
	specifyLanguage             string // 输出音频的语种. 参考链接：https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/language-support?tabs=tts
	outputVoice                 string // 输出音频的语种+口音. 参考链接：同上链接
	speechKey                   string
	speechRegion                string
	SpeechSynthesisOutputFormat common.SpeechSynthesisOutputFormat // 输出音频格式 默认是：common.Riff16Khz16BitMonoPcm
}

type LanguageCheckMode int

const (
	AutoCheck LanguageCheckMode = 0 // 自动检测模式
	Specify   LanguageCheckMode = 1 // 指定语言模式
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
	sid               int64
	speechConfig      *speech.SpeechConfig
	speechSynthesizer *speech.SpeechSynthesizer
}

func NewTTS(sid int64, con int, cfg *Config) (*TTS, error) {
	tts := &TTS{
		sid: sid,
	}
	tts.HttpSender = ttscommon.NewHttpSender(sid, con, tts.streamAsk)
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

// streamAsk 将文本发送到tts
func (tts *TTS) streamAsk(ctx context.Context, text string) (io.ReadCloser, error) {
	task := tts.speechSynthesizer.StartSpeakingTextAsync(text)

	var outcome speech.SpeechSynthesisOutcome

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("等待tts返回时被打断")
	case outcome = <-task: // 返回服务端的语音识别结果基本信息（此时还不包含音频数据）
	}

	defer outcome.Close()
	if outcome.Error != nil {
		return nil, fmt.Errorf("tts服务端对text合成语音过程中报错：%v", outcome.Error)
	}
	stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
	if err != nil {
		return nil, fmt.Errorf("获取tts语音合成结果的stream对象失败：%v", err)
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

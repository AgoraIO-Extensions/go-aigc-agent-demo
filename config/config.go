package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
	"time"
)

/* --------------------------------------------------  log  --------------------------------------------------------- */

type logConfig struct {
	File  string `toml:"file"`
	Level string `toml:"level"` // 只有指定了file值才会生效
}

/* --------------------------------------------------  rtc  --------------------------------------------------------- */

type rtc struct {
	AppID         string `toml:"app_id"`
	ChannelName   string `toml:"channel_name"`
	UserID        string `toml:"user_id"`
	Region        string `toml:"region"`
	OpenMsgReturn bool   `toml:"open_msg_return"`
}

/* -------------------------------------------------  filer  -------------------------------------------------------- */

type Filter struct {
	Open         bool `toml:"open"`
	BeginWinSize int  `toml:"begin_win_size"`
	EndWinSize   int  `toml:"end_win_size"`
	DataWinSize  int  `toml:"data_win_size"` // 包含了 beginWinSize
}

/* ------------------------------------------------  stt/tts  ------------------------------------------------------- */

type msSTT struct {
	SpeechKey              string   `toml:"speech_key"`
	SpeechRegion           string   `toml:"speech_region"`
	LanguageCheckMode      int      `toml:"language_check_mode"`       // 语言检测模式
	AutoAudioCheckLanguage []string `toml:"auto_audio_check_language"` // 自动识别音频语种范围
	SpecifyLanguage        string   `toml:"specify_language"`          // 指定的音频识别语种
	SetLog                 bool     `toml:"set_log"`                   // 记录speechsdk中的stt执行的信息日志
	ResultQueueSize        int      `toml:"result_queue_size"`         // stt返回文本的缓冲队列长度
}

type aliSTT struct {
	URL                string `toml:"url"`
	AKID               string `toml:"akid"`
	AKKey              string `toml:"akkey"`
	AppKey             string `toml:"appkey"`
	MaxSentenceSilence int    `toml:"maxSentenceSilence"`
}

type STTMode string

type SttSelect string
type TTSSelect string

const (
	MsSTT      SttSelect = "ms"
	AliSTT     SttSelect = "ali"
	MsTTS      TTSSelect = "ms"
	AliTTS     TTSSelect = "ali"
	AliCosyTTS TTSSelect = "cosy"
)

type STT struct {
	Select SttSelect `toml:"select"`
	Mode   STTMode   `toml:"mode"`
	MS     msSTT     `toml:"ms"`
	Ali    aliSTT    `toml:"ali"`
}

type SpeechSynthesisOutputFormat string

type msTTS struct {
	SpeechKey                   string                      `toml:"speech_key"`
	SpeechRegion                string                      `toml:"speech_region"`
	SetLog                      bool                        `toml:"set_log"`
	LanguageCheckMode           int                         `toml:"language_check_mode"`
	SpecifyLanguage             string                      `toml:"specify_language"`               // 输出音频的语种. 参考链接：https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/language-support?tabs=tts
	OutputVoice                 string                      `toml:"output_voice"`                   // 输出音频的语种+口音. 参考链接：同上链接
	SpeechSynthesisOutputFormat SpeechSynthesisOutputFormat `toml:"speech_synthesis_output_format"` // 输出音频格式
}

type aliTTS struct {
	URL    string `toml:"url"`
	AKID   string `toml:"akid"`
	AKKey  string `toml:"akkey"`
	AppKey string `toml:"appkey"`
}

type cosyTTS struct {
	URL    string `toml:"url"`
	AKID   string `toml:"akid"`
	AKKey  string `toml:"akkey"`
	AppKey string `toml:"appkey"`
}

type TTS struct {
	Select TTSSelect `toml:"select"`
	MS     msTTS     `toml:"ms"`
	Ali    aliTTS    `toml:"ali"`
	Cosy   cosyTTS   `toml:"cosy"`
}

/* --------------------------------------------------  llm  --------------------------------------------------------- */

type ModelSelect string

const (
	LLMQwen      ModelSelect = "qwen"
	LLMChatGPT4o ModelSelect = "chat-gpt4o"
)

type Prompt struct {
	OutputLanguage []string `toml:"output_language"`
	Prompt         string   `toml:"prompt"`
}

func (p *Prompt) Generate() string {
	prompt := p.Prompt
	if len(p.OutputLanguage) != 0 {
		prompt = prompt + fmt.Sprintf("请使用%s回答我。", p.OutputLanguage[0])
	}
	return prompt
}

type ClauseMode string

const (
	NoClause          ClauseMode = "none"
	PunctuationClause ClauseMode = "punctuation"
)

type QWen struct {
	Model      string `toml:"model"`
	URL        string `toml:"url"`
	DialogNums int    `toml:"dialog_nums"`
	ApiKey     string `toml:"apikey"`
}

type ChatGPT struct {
	Key        string `toml:"key"`
	Model      string `toml:"model"`
	EndPoint   string `toml:"end_point"`
	DialogNums int    `toml:"dialog_nums"`
}

type LLM struct {
	ModelSelect ModelSelect `toml:"model_select"` // 枚举值没定义自定义类型是因为flag解析命令行参数时无法使用自定义类型
	WithHistory bool        `toml:"with_history"`
	ClauseMode  ClauseMode  `toml:"clause_mode"`
	Prompt      Prompt      `toml:"prompt"`
	QWen        QWen        `toml:"qwen"`
	ChatGPT4o   ChatGPT     `toml:"chat_gpt4o"`
}

/* -------------------------------------------------  config  ------------------------------------------------------- */

var config *Config

type Config struct {
	StartTime   int64     // aigc-worker任务启动时间(unix s)
	MaxLifeTime int64     `toml:"max_life_time"` // 最大生命周期（s）
	RTC         rtc       `toml:"rtc"`
	Filter      Filter    `toml:"filter"`
	STT         STT       `toml:"stt"`
	TTS         TTS       `toml:"tts"`
	Log         logConfig `toml:"log"`
	LLM         LLM       `toml:"llm"`
}

func Inst() *Config {
	return config
}

// Init 基于配置文件初始化Config
func Init(filePath string) error {
	config = &Config{
		StartTime: time.Now().Unix(),
	}
	if filePath != "" {
		if err := config.load(filePath); err != nil {
			return err
		}
	}

	return nil
}

// The ErrConfigValidationFailed error is used so that external callers can do a type assertion
// to defer handling of this specific error when someone does not want strict type checking.
// This is needed only because logging hasn't been set up at the time we parse the config file.
// This should all be ripped out once strict config checking is made the default behavior.
type ErrConfigValidationFailed struct {
	confFile       string
	UndecodedItems []string
}

func (e *ErrConfigValidationFailed) Error() string {
	return fmt.Sprintf("config file %s contained invalid configuration options: %s",
		e.confFile, strings.Join(e.UndecodedItems, ", "))
}

func (cfg *Config) load(confFile string) error {
	metaData, err := toml.DecodeFile(confFile, cfg)
	// If any items in confFile file are not mapped into the Config struct, issue
	// an error and stop the server from starting.
	undecoded := metaData.Undecoded()
	if len(undecoded) > 0 && err == nil {
		var undecodedItems []string
		for _, item := range undecoded {
			undecodedItems = append(undecodedItems, item.String())
		}
		err = &ErrConfigValidationFailed{confFile, undecodedItems}
	}

	return err
}

package ali

import (
	"fmt"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"log"
	"os"
)

type Config struct {
	URL    string
	AppKey string
	Token  string

	nlsStartParamsConf nls.SpeechTranscriptionStartParam
	nlsConnConf        *nls.ConnectionConfig
	nlsLogger          *nls.NlsLogger

	connPool *connPool
}

func Init(url, appkey, token string) (*Config, error) {
	cfg := &Config{
		URL:    url,
		AppKey: appkey,
		Token:  token,
	}
	cfg.nlsStartParamsConf = nls.SpeechTranscriptionStartParam{
		Format:                         "pcm",
		SampleRate:                     16000,
		EnableIntermediateResult:       true,
		EnablePunctuationPrediction:    true,
		EnableInverseTextNormalization: true,
		MaxSentenceSilence:             3000,
		EnableWords:                    false,
	}
	cfg.nlsConnConf = nls.NewConnectionConfigWithToken(url, appkey, token)
	cfg.nlsLogger = nls.NewNlsLogger(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
	cfg.nlsLogger.SetLogSil(true)
	cfg.nlsLogger.SetDebug(false)

	pool, err := initConnPool(5, cfg)
	if err != nil {
		return nil, fmt.Errorf("[initConnPool]Failed to initialize the connection pool for ali-stt.%v", err)
	}
	cfg.connPool = pool

	return cfg, nil
}

package sentence

import (
	"context"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

var FirstSid = time.Now().Unix() * 10

type Stage int

const (
	BeforeSendToRTC Stage = 0
	OnSendToRTC     Stage = 1
	//AfterSendToRTC  Stage = 2
)

type MetaData struct {
	Sid                    int64
	Sgid                   int64
	Stage                  Stage
	FilterAudioTailRcvTime time.Time // the time of the filter outputting the tail chunk
}

// LogHook set the logs to automatically print sid and sgid
func LogHook(ctx context.Context, record *slog.Record) {
	if sMetaData, ok := ctx.Value(logger.SentenceMetaData).(*MetaData); ok {
		record.AddAttrs(slog.Int64("sid", sMetaData.Sid))
		record.AddAttrs(slog.Int64("sgid", sMetaData.Sgid))
	}
}

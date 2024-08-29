package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func Init(file string, level string, uuid string) error {
	if file == "" {
		logger, _ = zap.NewProduction()
	} else {
		wcore := zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    40,   // megabytes
			MaxAge:     31,   // days
			MaxBackups: 10,   // the maximum number of old log files to retain
			Compress:   true, // use gzip to compress all rotated log files
		})
		zaplevel := zap.NewAtomicLevel()
		zaplevel.UnmarshalText([]byte(level))
		zapencCfg := zap.NewProductionEncoderConfig()
		zapencCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		zapencCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			loc := time.FixedZone("UTC+8", 8*60*60)
			enc.AppendString(t.In(loc).Format("2006-01-02T15:04:05.000Z07:00"))
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zapencCfg),
			wcore,
			zaplevel.Level(),
		)

		// append "caller", "pid" to log
		opts := []zap.Option{
			zap.AddCaller(),
			zap.Fields(zap.Any("uuid", uuid)),
		}
		logger = zap.New(core, opts...)
	}

	return nil
}

func AddOptions(opts ...zap.Option) {
	logger = logger.WithOptions(opts...)
}

func Inst() *zap.Logger {
	return logger
}

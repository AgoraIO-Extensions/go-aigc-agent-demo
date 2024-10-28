package ms

import (
	"context"
	"go-aigc-agent-demo/pkg/logger"
	"io"
	"log/slog"
	"sync"
)

func PreConn(con int) {
	ctx := context.Background()
	tts := &TTS{ctx: ctx}

	wg := sync.WaitGroup{}
	for i := 0; i < con; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			rc, err := tts.streamAsk(ctx, -1, "hello")
			if err != nil {
				logger.ErrorContext(ctx, "[tts] establish pre-connection failed", slog.Int("i", i), slog.Any("err", err))
				return
			}
			if _, err = io.ReadAll(rc); err != nil {
				logger.ErrorContext(ctx, "[tts][io.ReadAll]", slog.Int("i", i), slog.Any("err", err))
				return
			}
			rc.Close()
			logger.InfoContext(ctx, "[tts] establish pre-connection succeeded", slog.Int("i", i))
		}(i)
	}

	wg.Wait()
}

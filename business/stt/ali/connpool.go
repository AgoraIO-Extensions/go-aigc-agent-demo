package ali

import (
	"context"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

type connPool struct {
	cfg      *Config
	pool     chan *conn
	poolSize int // Maximum number of idle connections in the connection pool, recommended value: 5
}

func initConnPool(size int, cfg *Config) (*connPool, error) {
	p := &connPool{
		cfg:      cfg,
		pool:     make(chan *conn, size),
		poolSize: size,
	}
	wg := new(sync.WaitGroup)
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c, err := newConn(cfg, context.Background())
			if err != nil {
				logger.Error("[stt]Failed to initialize the connection pool for ali-stt", slog.Any("err", err), slog.Int("i", i))
			}
			p.pool <- c
		}(i)
	}
	wg.Wait()
	p.generateConn()
	return p, nil
}

func (p *connPool) generateConn() {
	go func() {
		ticker := time.Tick(time.Second)

		for range ticker {
			n := len(p.pool)
			if n < 2 {
				logger.Info("[stt] ali-stt Number of connections in the connection pool<2", slog.Int("conn_num", n))
			}
			if n == p.poolSize {
				<-p.pool
			}

			concurrency := 1
			if n <= p.poolSize-2 {
				concurrency = 2
			}
			for i := 0; i < concurrency; i++ {
				go func() {
					start := time.Now()
					c, err := newConn(p.cfg, context.Background())
					if err != nil {
						logger.Error("[stt]Failed to asynchronously create connection for ali-stt", slog.Any("err", err))
						return
					}
					dur := time.Now().Sub(start)
					if dur > time.Second {
						logger.Info("[stt]ali-stt Connection establishment time>1s", slog.Int64("dur", dur.Milliseconds()))
					}
					p.pool <- c
				}()
			}
		}
	}()
}

func (p *connPool) GetConn() *conn {
	var c *conn
	for c == nil {
		c = <-p.pool
		if time.Now().After(c.expTime) {
			c = nil
		}
	}
	return c
}

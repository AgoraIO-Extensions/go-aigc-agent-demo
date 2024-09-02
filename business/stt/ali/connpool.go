package ali

import (
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

type connPool struct {
	cfg      *Config
	pool     chan *conn
	poolSize int // 连接池最大空闲连接数，建议取值5
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
			c, err := newConn(cfg, 0)
			if err != nil {
				logger.Inst().Error("[stt]ali-stt初始化连接池失败", zap.Error(err), zap.Int("i", i))
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
				logger.Inst().Info("[stt] ali-stt 连接池中的连接数<2", zap.Int("conn_num", n))
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
					c, err := newConn(p.cfg, 0)
					if err != nil {
						logger.Inst().Error("[stt]ali-stt异步生成连接失败", zap.Error(err))
						return
					}
					dur := time.Now().Sub(start)
					if dur > time.Second {
						logger.Inst().Info("[stt]ali-stt 连接池建立连接耗时>1s", zap.Int64("dur", dur.Milliseconds()))
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

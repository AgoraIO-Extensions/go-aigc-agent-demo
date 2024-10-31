package util

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

type BufPipe struct {
	buf    *bytes.Buffer
	mu     *sync.Mutex
	cond   *sync.Cond
	closed bool
}

func NewBufPipe() *BufPipe {
	bp := &BufPipe{
		buf: new(bytes.Buffer),
		mu:  new(sync.Mutex),
	}
	bp.cond = sync.NewCond(bp.mu)
	return bp
}

func (bp *BufPipe) Write(p []byte) (n int, err error) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if bp.closed {
		return 0, errors.New("can't write to a closed BufPipe")
	}

	n, err = bp.buf.Write(p)
	bp.cond.Signal()
	return n, err
}

func (bp *BufPipe) Read(p []byte) (n int, err error) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	for bp.buf.Len() == 0 {
		if bp.closed {
			return 0, io.EOF
		}
		bp.cond.Wait()
	}
	return bp.buf.Read(p)
}

func (bp *BufPipe) Close() error {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.closed = true
	bp.cond.Broadcast()
	return nil
}

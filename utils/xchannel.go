package utils

import (
	"errors"
	"sync"
)

type XChannel struct {
	rw    sync.RWMutex
	close bool
	Ch    chan []byte
}

func NewXChannel(size int) *XChannel {
	return &XChannel{
		Ch: make(chan []byte, size),
	}
}

func (x *XChannel) Send(msg []byte) error {
	x.rw.RLock()
	defer x.rw.RUnlock()

	if x.close {
		return errors.New("channel close")
	}

	x.Ch <- msg
	return nil
}

func (x *XChannel) Close() {
	x.rw.Lock()
	defer x.rw.Unlock()

	if x.close == false {
		x.close = true
		close(x.Ch)
	}
}

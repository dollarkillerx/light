package utils

import (
	"errors"
	"sync"

	"github.com/dollarkillerx/light/protocol"
)

type XChannel struct {
	rw    sync.RWMutex
	close bool
	ch    chan *protocol.Message
}

func NewXChannel(size int) *XChannel {
	return &XChannel{
		ch: make(chan *protocol.Message, size),
	}
}

func (x *XChannel) Send(msg *protocol.Message) error {
	x.rw.RLock()
	defer x.rw.Unlock()

	if x.close {
		return errors.New("channel close")
	}
	x.ch <- msg

	return nil
}

func (x *XChannel) Read() <-chan *protocol.Message {
	return x.ch
}

func (x *XChannel) Close() {
	x.rw.Lock()
	defer x.rw.Unlock()

	x.close = true
	close(x.ch)
}

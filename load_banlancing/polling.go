package load_banlancing

import (
	"sync"

	"github.com/dollarkillerx/light/discovery"
)

// Polling 轮训
type Polling struct {
	mu  sync.Mutex
	ser []*discovery.Server
	idx int
}

func NewPolling() *Polling {
	return &Polling{
		ser: []*discovery.Server{},
		idx: 0,
	}
}

func (p *Polling) InitBalancing(sers []*discovery.Server) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ser = sers
}

func (p *Polling) GetService() *discovery.Server {
	if p.idx >= len(p.ser)-1 {
		p.idx = 0
	}

	return p.ser[p.idx]
}

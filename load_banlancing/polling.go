package load_banlancing

import (
	"github.com/pkg/errors"
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
	p.idx = 0
}

func (p *Polling) GetService() (*discovery.Server, error) {
	if p.idx >= len(p.ser)-1 {
		p.idx = 0
	}

	if len(p.ser) == 0 {
		return nil, errors.New("404 Discovery")
	}
	return p.ser[p.idx], nil
}

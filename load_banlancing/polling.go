package load_banlancing

import "github.com/dollarkillerx/light/discovery"

// Polling 轮训
type Polling struct {
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
	p.ser = sers
}

func (p *Polling) GetService() *discovery.Server {
	if p.idx >= len(p.ser)-1 {
		p.idx = 0
	}

	return p.ser[p.idx]
}

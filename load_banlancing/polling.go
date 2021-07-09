package load_banlancing

// Polling 轮训
type Polling struct {
	ser []string
	idx int
}

func NewPolling() *Polling {
	return &Polling{
		ser: []string{},
		idx: 0,
	}
}

func (p *Polling) InitBalancing(sers []string) {
	p.ser = sers
}

func (p *Polling) GetService() string {
	if p.idx >= len(p.ser)-1 {
		p.idx = 0
	}

	return p.ser[p.idx]
}

package currencies

import (
	"sync"
	"time"
)

type Source string

const (
	EXCHANGE Source = "exchange"
	INTERNAL Source = "internal"
)

type Pricer struct {
	sync.RWMutex

	Using   Source
	Sources map[Source]*Price
}

func (p *Pricer) SetSource(s Source) {
	p.Lock()
	defer p.Unlock()

	p.Using = s
}

func (p *Pricer) GetSource() Source {
	p.RLock()
	defer p.RUnlock()

	return p.Using
}

func (p *Pricer) GetPrice() (uint64, Source, time.Time) {
	p.RLock()
	defer p.RUnlock()

	if p.Sources[p.Using] == nil {
		return 0, "", time.Now()
	}

	price, updated := p.Sources[p.Using].Get()
	return price, p.Using, updated
}

func (p *Pricer) SetPrice(s Source, a uint64) {
	p.Lock()
	defer p.Unlock()

	if p.Sources[s] == nil {
		p.Sources[s] = NewPrice(a)
		return
	}

	p.Sources[s].Set(a)
}

package currencies

import (
	"sync"
	"time"
)

type Price struct {
	sync.RWMutex

	Updated time.Time
	Amount  uint64
}

func NewPrice(amount uint64) *Price {
	return &Price{
		Updated: time.Now(),
		Amount:  amount,
	}
}

func (p *Price) Get() (uint64, time.Time) {
	p.RLock()
	defer p.RUnlock()

	return p.Amount, p.Updated
}

func (p *Price) Set(amount uint64) {
	p.Lock()
	defer p.Unlock()

	p.Amount = amount
	p.Updated = time.Now()
}

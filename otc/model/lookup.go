package model

import (
	"sync"

	"github.com/skycoin/services/otc/types"
)

// Lookup is used to manage Drop/Address associations for navigating Model disk
// storage.
type Lookup struct {
	sync.RWMutex

	dropToSky map[types.Currency]map[types.Drop]types.Address
}

func NewLookup() *Lookup {
	return &Lookup{
		dropToSky: make(map[types.Currency]map[types.Drop]types.Address, 0),
	}
}

// GetAddress returns the Address associated with a Drop. An error is returned
// if no Address exists for the Drop.
func (l *Lookup) GetAddress(d types.Drop, c types.Currency) (types.Address, error) {
	l.RLock()
	defer l.RUnlock()

	if l.dropToSky[c] == nil || l.dropToSky[c][d] == "" {
		return "", ErrDropMissing
	}
	return l.dropToSky[c][d], nil
}

// SetDrop associates a Drop with an Address (SKY) in the Lookup map.
func (l *Lookup) SetDrop(d types.Drop, c types.Currency, a types.Address) {
	l.Lock()
	defer l.Unlock()

	if l.dropToSky[c] == nil {
		l.dropToSky[c] = make(map[types.Drop]types.Address, 0)
	}
	l.dropToSky[c][d] = a
}

package dropper

import (
	"errors"
	"sync"
	"time"

	"github.com/skycoin/services/otc/types"
)

type Source string

const (
	EXCHANGE Source = "exchange"
	INTERNAL Source = "internal"
)

type Value struct {
	Updated time.Time
	Amount  uint64
}

type Dropper struct {
	Connections types.Connections

	ValueMutex   sync.RWMutex
	ValueSource  Source
	Value        map[types.Currency]uint64
	ValueUpdated map[types.Currency]time.Time
}

func NewDropper(config *types.Config) (*Dropper, error) {
	btc, err := NewBTCConnection(config)

	return &Dropper{
		Connections:  types.Connections{types.BTC: btc},
		ValueSource:  EXCHANGE,
		Value:        make(map[types.Currency]uint64, 0),
		ValueUpdated: make(map[types.Currency]time.Time, 0),
	}, err
}

var ErrConnectionMissing = errors.New("connection doesn't exist")

func (d *Dropper) GetBalance(c types.Currency, a types.Drop) (uint64, error) {
	connection, exists := d.Connections[c]
	if !exists {
		return 0.0, ErrConnectionMissing
	}

	return connection.Balance(a)
}

// GetValueSource gets the reference used to determine currency value.
func (d *Dropper) GetValueSource() Source {
	d.ValueMutex.RLock()
	defer d.ValueMutex.RUnlock()

	return d.ValueSource
}

// SetValueSource sets the reference for determining currency value. Currently
// only EXCHANGE and INTERNAL are supported.
func (d *Dropper) SetValueSource(s Source) {
	d.ValueMutex.Lock()
	defer d.ValueMutex.Unlock()

	d.ValueSource = s
}

// SetValue sets the value of 1 SKY (amount) for the currency.
func (d *Dropper) SetValue(c types.Currency, amount uint64) {
	d.ValueMutex.Lock()
	defer d.ValueMutex.Unlock()

	d.ValueUpdated[c] = time.Now().In(time.UTC)
	d.Value[c] = amount
}

// GetValue returns the equivalent of 1 SKY in the passed currency.
func (d *Dropper) GetValue(c types.Currency) (uint64, error) {
	d.ValueMutex.RLock()
	defer d.ValueMutex.RUnlock()

	if d.ValueSource == EXCHANGE {
		if value, err := d.Connections[c].Value(); err != nil {
			return 0, err
		} else {
			d.ValueUpdated[c] = time.Now().UTC()
			return value, nil
		}
	} else if d.ValueSource == INTERNAL {
		if _, exists := d.Value[c]; !exists {
			return 0, ErrConnectionMissing
		} else {
			return d.Value[c], nil
		}
	}

	// should never be reached
	return 0, nil
}

func (d *Dropper) GetUpdated(c types.Currency) time.Time {
	d.ValueMutex.RLock()
	defer d.ValueMutex.RUnlock()

	return d.ValueUpdated[c]
}

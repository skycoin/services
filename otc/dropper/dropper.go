package dropper

import (
	"errors"
	"sync"
	"time"

	"github.com/skycoin/services/otc/exchange"
	"github.com/skycoin/services/otc/types"
)

type Source string

const (
	EXCHANGE Source = "exchange"
	INTERNAL Source = "internal"
)

type Dropper struct {
	Connections types.Connections
	Currencies  map[types.Currency]*Values
}

type Values struct {
	Using   Source
	Sources map[Source]*Value
}

func NewValues(s Source, v *Value) *Values {
	values := &Values{
		Using:   s,
		Sources: make(map[Source]*Value, 0),
	}

	// add initial
	values.Sources[s] = v

	return values
}

func (v *Values) SetSource(s Source) { v.Using = s }

func (v *Values) GetSource() Source { return v.Using }

func (v *Values) GetValue() (uint64, time.Time) {
	if v.Sources[v.Using] == nil {
		return 0, time.Now()
	}

	return v.Sources[v.Using].Get()
}

func (v *Values) SetValue(s Source, a uint64) {
	if v.Sources[s] == nil {
		v.Sources[s] = NewValue(a)
		return
	}

	v.Sources[s].Set(a)
}

type Value struct {
	sync.RWMutex
	Updated time.Time
	Amount  uint64
}

func NewValue(a uint64) *Value {
	return &Value{
		Updated: time.Now(),
		Amount:  a,
	}
}

func (v *Value) Get() (uint64, time.Time) {
	v.RLock()
	defer v.RUnlock()
	return v.Amount, v.Updated
}

func (v *Value) Set(a uint64) {
	v.Lock()
	defer v.Unlock()
	v.Amount = a
	v.Updated = time.Now()
}

func NewDropper(config *types.Config) (*Dropper, error) {
	btc, err := NewBTCConnection(config)
	if err != nil {
		return nil, err
	}

	d := &Dropper{
		Connections: types.Connections{types.BTC: btc},
		Currencies:  make(map[types.Currency]*Values, 0),
	}

	d.Currencies[types.BTC] = NewValues(
		// default price source
		INTERNAL,
		// get default price from config
		NewValue(config.Dropper.BTC.Price),
	)

	return d, nil
}

var ErrConnectionMissing = errors.New("connection doesn't exist")

func (d *Dropper) GetBalance(c types.Currency, a types.Drop) (uint64, error) {
	connection, exists := d.Connections[c]
	if !exists {
		return 0.0, ErrConnectionMissing
	}

	return connection.Balance(a)
}

// TODO: support more currencies / configs
func (d *Dropper) Start() {
	go func() {
		for {
			val, err := exchange.GetBTCValue()
			if err != nil {
				// if error with exchange, set new source to internal
				//
				// TODO: check if this is a good idea
				d.Currencies[types.BTC].SetSource(INTERNAL)
			}

			d.Currencies[types.BTC].SetValue(EXCHANGE, val)

			<-time.After(time.Second * 15)
		}
	}()
}

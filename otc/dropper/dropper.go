package dropper

import (
	"errors"

	"github.com/skycoin/services/otc/types"
)

type Dropper struct {
	Connections types.Connections
}

func NewDropper(config *types.Config) (*Dropper, error) {
	btc, err := NewBTCConnection(config)

	return &Dropper{
		Connections: types.Connections{types.BTC: btc},
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

// GetValue returns SKY value of the amount of currency.
func (d *Dropper) GetValue(c types.Currency, amount uint64) uint64 {
	// use exchange api or value set by admin panel
	return 0
}

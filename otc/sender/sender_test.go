package sender

import (
	"flag"
	"testing"

	"github.com/skycoin/services/otc/dropper"
	"github.com/skycoin/services/otc/skycoin"
)

const (
	SKYCOIN_NODE = "localhost:6430"
)

var (
	SKYCOIN *skycoin.Connection
	DROPPER *dropper.Dropper

	SKYCOIN_SEED = flag.String(
		"skycoin_seed",
		"",
		"seed used for skycoin wallet daemon (sending coins)",
	)
)

func init() {
	flag.Parse()

	var err error

	SKYCOIN, err = skycoin.NewConnection(SKYCOIN_NODE, *SKYCOIN_SEED)
	if err != nil {
		panic(err)
	}

	DROPPER, err = dropper.NewDropper()
	if err != nil {
		panic(err)
	}
}

func TestNewSender(t *testing.T) {
	s, err := NewSender(SKYCOIN, DROPPER)
	if err != nil {
		panic(err)
	}

	if s == nil {
		t.Fatal("nil sender")
	}
}

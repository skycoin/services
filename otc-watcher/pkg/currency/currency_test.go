package currency

import (
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

type TestMock struct{}

func (t *TestMock) Stop() error                            { return nil }
func (t *TestMock) Scan(h uint64) (chan *otc.Block, error) { return nil, nil }
func (t *TestMock) Get(h uint64) (*otc.Block, error)       { return nil, nil }
func (t *TestMock) Height() (uint64, error)                { return 0, nil }

func TestGet(t *testing.T) {
	cons := Connections(map[otc.Currency]Connection{otc.BTC: &TestMock{}})

	if _, err := cons.Get(otc.SKY, 0); err != ErrConnMissing {
		t.Fatal("connection should be missing")
	}

	if _, err := cons.Get(otc.BTC, 0); err != nil {
		t.Fatal(err)
	}
}

func TestHeight(t *testing.T) {
	cons := Connections(map[otc.Currency]Connection{otc.BTC: &TestMock{}})

	if _, err := cons.Height(otc.SKY); err != ErrConnMissing {
		t.Fatal("connection should be missing")
	}

	if _, err := cons.Height(otc.BTC); err != nil {
		t.Fatal(err)
	}
}

func TestScan(t *testing.T) {
	cons := Connections(map[otc.Currency]Connection{otc.BTC: &TestMock{}})

	if _, err := cons.Scan(otc.SKY, 0); err != ErrConnMissing {
		t.Fatal("connection should be missing")
	}

	if _, err := cons.Scan(otc.BTC, 0); err != nil {
		t.Fatal(err)
	}
}

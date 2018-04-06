package scanner

import (
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

func TestStorageNew(t *testing.T) {
	storage := NewStorage(otc.BTC)

	if storage.Filename != string(otc.BTC)+".json" {
		t.Fatal("bad filename")
	}

	if storage.Addresses == nil {
		t.Fatal("didn't initialize addresses")
	}
}

func TestStorageRegister(t *testing.T) {
	storage := NewStorage(otc.BTC)
	storage.Register("address")

	if storage.Addresses["address"] == nil {
		t.Fatal("didn't register address")
	}
}

func TestStorageOutputs(t *testing.T) {
	storage := NewStorage(otc.BTC)
	storage.Register("address")
	storage.Addresses["address"].Outputs.Update(
		"hash",
		1,
		&otc.OutputVerbose{
			Amount:        32,
			Confirmations: 1,
			Height:        40,
		},
	)
	outputs := storage.Outputs("address")

	if outputs["hash"][1].Amount != 32 {
		t.Fatal("outputs failed")
	}
}

func TestStorageUpdate(t *testing.T) {
	storage := NewStorage(otc.BTC)
	storage.Register("address")

	storage.Update(&otc.Block{
		Height: 32,
		Transactions: map[string]*otc.Transaction{
			"transaction": &otc.Transaction{
				Hash:          "transaction",
				Confirmations: 3,
				Out: map[int]*otc.Output{
					1: &otc.Output{
						Amount:    32000000,
						Addresses: []string{"address"},
					},
					2: &otc.Output{
						Amount:    100000000,
						Addresses: []string{"irrelevant"},
					},
					3: &otc.Output{
						Amount:    50000000,
						Addresses: []string{"address"},
					},
				},
			},
		},
	})

	if storage.Addresses["address"].Outputs["transaction"][1] == nil ||
		storage.Addresses["address"].Outputs["transaction"][3] == nil ||
		len(storage.Addresses["address"].Outputs["transaction"]) != 2 {
		t.Fatal("update failed")
	}
}

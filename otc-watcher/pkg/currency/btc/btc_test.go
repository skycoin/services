package btc

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Mock struct {
	getBlockHash      func(int64) (*chainhash.Hash, error)
	getBlockVerboseTx func(*chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	getBlockCount     func() (int64, error)
	waitForShutdown   func()
}

func (m *Mock) GetBlockHash(h int64) (*chainhash.Hash, error) {
	return m.getBlockHash(h)
}

func (m *Mock) GetBlockVerboseTx(h *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	return m.getBlockVerboseTx(h)
}

func (m *Mock) GetBlockCount() (int64, error) {
	return m.getBlockCount()
}

func (m *Mock) WaitForShutdown() {
	m.waitForShutdown()
}

///////////////////////////////////////////////////////////////////////////////

func GetBlockHash(err error) func(int64) (*chainhash.Hash, error) {
	return func(h int64) (*chainhash.Hash, error) {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func GetBlockVerboseTx(err error, a float64) func(*chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	return func(h *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
		if err != nil {
			return nil, err
		}

		return &btcjson.GetBlockVerboseResult{
			RawTx: []btcjson.TxRawResult{
				{
					Hash:          "hash",
					Confirmations: 3,
					Vout: []btcjson.Vout{
						{
							Value: a,
							N:     1,
							ScriptPubKey: btcjson.ScriptPubKeyResult{
								Addresses: []string{"address"},
							},
						},
					},
				},
			},
		}, nil
	}
}

func GetBlockCount(c int64, err error) func() (int64, error) {
	return func() (int64, error) {
		return c, err
	}
}

func WaitForShutdown(s chan struct{}) func() {
	return func() {
		s <- struct{}{}
	}
}

///////////////////////////////////////////////////////////////////////////////

func TestStop(t *testing.T) {
	var (
		stop = make(chan struct{}, 1)
		shut = make(chan struct{}, 1)
	)

	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, 0),
			GetBlockCount(0, nil),
			WaitForShutdown(shut),
		},
		Account: "",
		stop:    stop,
	}

	if err := connection.Stop(); err != nil {
		t.Fatal(err)
	}

	select {
	case <-stop:
		break
	default:
		t.Fatal("stop chan not filled")
	}

	select {
	case <-shut:
		break
	default:
		t.Fatal("client shutdown not called")
	}
}

func TestHeight(t *testing.T) {
	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, 0),
			GetBlockCount(32, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if height, err := connection.Height(); height != 32 || err != nil {
		t.Fatal("couldn't get height")
	}
}

func TestGetGood(t *testing.T) {
	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, 0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Get(32); err != nil {
		t.Fatal("couldn't get block")
	}
}

func TestGetBadHash(t *testing.T) {
	bad := fmt.Errorf("bad error!")

	connection := &Connection{
		Client: &Mock{
			GetBlockHash(bad),
			GetBlockVerboseTx(nil, 0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Get(32); err != bad {
		t.Fatal("should have returned an error")
	}
}

func TestGetBadVerboseTx(t *testing.T) {
	bad := fmt.Errorf("bad error!")

	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(bad, 1.0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Get(32); err != bad {
		t.Fatal("should have returned an error")
	}
}

func TestGetBadAmount(t *testing.T) {
	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, math.NaN()),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Get(32); err == nil {
		t.Fatal("should have returned an error")
	}
}

func TestScanStop(t *testing.T) {
	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, 1.0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    make(chan struct{}, 1),
	}

	connection.stop <- struct{}{}

	if _, err := connection.Scan(32); err != nil {
		t.Fatal(err)
	}
}

func TestScanMaxed(t *testing.T) {
	var buf bytes.Buffer

	connection := &Connection{
		Logs: log.New(&buf, "", 0),
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(fmt.Errorf(
				"-1: Block number out of range",
			), 1.0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Scan(32); err != nil {
		t.Fatal(err)
	}

	// wait for goroutine to log error
	<-time.After(time.Second / 10)

	if buf.String() != "waiting for block: 32\n" {
		t.Fatalf("expected error log, got '%s'\n", buf.String())
	}
}

func TestScanErr(t *testing.T) {
	var buf bytes.Buffer

	connection := &Connection{
		Logs: log.New(&buf, "", 0),
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(fmt.Errorf("bad"), 1.0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	if _, err := connection.Scan(32); err != nil {
		t.Fatal(err)
	}

	// wait for goroutine to log error
	<-time.After(time.Second / 10)

	if buf.String() != "scan error: bad\n" {
		t.Fatalf("expected error log, got '%s'\n", buf.String())
	}
}

func TestScan(t *testing.T) {
	connection := &Connection{
		Client: &Mock{
			GetBlockHash(nil),
			GetBlockVerboseTx(nil, 1.0),
			GetBlockCount(0, nil),
			WaitForShutdown(nil),
		},
		Account: "",
		stop:    nil,
	}

	blocks, err := connection.Scan(32)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second / 10)

	select {
	case <-blocks:
		break
	default:
		t.Fatal("blocks not send from scanner")
	}
}

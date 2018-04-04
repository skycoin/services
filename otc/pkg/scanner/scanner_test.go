package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/services/otc/pkg/watcher"
)

type MockClient struct {
	Do func(req *http.Request) (*http.Response, error)
}

func (c *MockClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.Do(req)
}

func MockWatcher(kind string) *watcher.Watcher {
	out := map[string]map[int]*otc.OutputVerbose{
		"transaction": {
			1: {
				Amount:        100000,
				Confirmations: 1,
				Addresses:     []string{"address"},
				Height:        500000,
			},
		},
	}

	return &watcher.Watcher{
		Client: &http.Client{
			Transport: &MockClient{
				func(req *http.Request) (*http.Response, error) {
					res := httptest.NewRecorder()

					if kind == "error" {
						return nil, fmt.Errorf("test error!")
					}

					// write mock outputs to result
					if err := json.NewEncoder(res).Encode(out); err != nil {
						return nil, err
					}

					return res.Result(), nil
				},
			},
		},
	}
}

func TestTaskGood(t *testing.T) {
	order, err := Task(MockWatcher(""))(&otc.User{
		Drop: &otc.Drop{
			Address:  "address",
			Currency: otc.BTC,
		},
	})

	if order == nil || err != nil {
		t.Fatal("bad scan")
	}
}

func TestTaskBad(t *testing.T) {
	order, err := Task(MockWatcher("error"))(&otc.User{
		Drop: &otc.Drop{
			Address:  "address",
			Currency: otc.BTC,
		},
	})

	if order != nil || err.Error() != "Post /outputs: test error!" {
		t.Fatalf(
			"expected 'Post /outputs: test error!', got '%s'\n",
			err.Error(),
		)
	}
}

func TestTaskExists(t *testing.T) {
	order, err := Task(MockWatcher(""))(&otc.User{
		Drop: &otc.Drop{
			Address:  "address",
			Currency: otc.BTC,
		},
		Orders: []*otc.Order{
			{
				Id: "transaction:1",
			},
		},
	})

	if order != nil || err != nil {
		t.Fatal("order should be empty")
	}
}

package watcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

func TestNew(t *testing.T) {
	conf := &otc.Config{}
	conf.Watcher.Node = "location"

	watcher, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}

	if watcher.Node != "location" {
		t.Fatal("bad node")
	}

	if watcher.Client == nil {
		t.Fatal("nil http client")
	}
}

type Mock struct {
	Type string
}

func (m *Mock) RoundTrip(req *http.Request) (*http.Response, error) {
	res := httptest.NewRecorder()

	if m.Type == "error" {
		return nil, fmt.Errorf("error!")
	} else if m.Type == "bad" {
		res.WriteHeader(500)
		return res.Result(), nil
	}

	outputs := map[string]map[int]*otc.OutputVerbose{
		"transaction": {
			1: {
				Amount: 1,
			},
		},
	}

	if err := json.NewEncoder(res).Encode(outputs); err != nil {
		return nil, err
	}

	return res.Result(), nil
}

func TestOutputs(t *testing.T) {
	watcher := &Watcher{
		Client: &http.Client{
			Transport: &Mock{},
		},
	}

	outputs, err := watcher.Outputs(&otc.Drop{"address", otc.BTC})
	if err != nil {
		t.Fatal(err)
	}

	if outputs == nil {
		t.Fatal("outputs shouldn't be nil")
	}
}

func TestOutputsError(t *testing.T) {
	watcher := &Watcher{
		Client: &http.Client{
			Transport: &Mock{"error"},
		},
	}

	if _, err := watcher.Outputs(&otc.Drop{"address", otc.BTC}); err == nil {
		t.Fatal("should be an error")
	}
}

func TestOutputsBad(t *testing.T) {
	watcher := &Watcher{
		Client: &http.Client{
			Transport: &Mock{"bad"},
		},
	}

	if _, err := watcher.Outputs(&otc.Drop{"address", otc.BTC}); err == nil {
		t.Fatal("should be an error")
	}
}

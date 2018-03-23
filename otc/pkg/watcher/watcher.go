package watcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/services/otc/pkg/otc"
)

type Watcher struct {
	Client *http.Client
	Node   string
}

func New(conf *otc.Config) (*Watcher, error) {
	return &Watcher{
		Client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 10,
		},
		Node: conf.Watcher.Node,
	}, nil
}

func (w *Watcher) Register(drop *otc.Drop) error {
	var buf bytes.Buffer

	// encode json
	if err := json.NewEncoder(&buf).Encode(drop); err != nil {
		return err
	}

	// send POST request to watcher
	resp, err := w.Client.Post(
		w.Node+"/register", "application/json", &buf,
	)
	if err != nil {
		return err
	}

	// check status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("watcher returned error")
	}

	return nil
}

func (w *Watcher) Outputs(drop *otc.Drop) (otc.Outputs, error) {
	var buf bytes.Buffer

	// encode json
	if err := json.NewEncoder(&buf).Encode(drop); err != nil {
		return nil, err
	}

	// send POST request to watcher
	reps, err := w.Client.Post(
		w.Node+"/outputs", "application/json", &buf,
	)
	if err != nil {
		return nil, err
	}

	// check status code
	if resp.StatusCode != 200 {
		return fmt.Errorf("watcher returned error")
	}

	// TODO: output unmarshalling

	return nil, nil
}

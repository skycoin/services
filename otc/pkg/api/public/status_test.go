package public

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func TestStatus(t *testing.T) {
	modl := &model.Model{
		Running: true,
		Logger:  log.New(ioutil.Discard, "", 0),
		Router:  actor.New(nil, nil),
		Lookup: map[string]*otc.Request{
			"currency:address": &otc.Request{
				Status: otc.NEW,
				Times: &otc.Times{
					UpdatedAt: 1,
				},
			},
		},
	}

	tests := [][]string{
		{
			`bad json`,
			`invalid JSON`,
		},
		{
			`{"drop_address":"bad","drop_currency":"BAD"}`,
			`request missing`,
		},
		{
			`{"drop_address":"address","drop_currency":"currency"}`,
			`{"status":"new","updated_at":1}`,
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		buf.WriteString(test[0])
		req := httptest.NewRequest("GET", "http:///", &buf)
		res := httptest.NewRecorder()

		Status(nil, modl)(res, req)

		out, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(out)) != test[1] {
			t.Fatalf(`expected "%s", got "%s"`, test[1],
				strings.TrimSpace(string(out)))
		}
	}
}

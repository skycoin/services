package admin

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/model"
)

func MockModel() *model.Model {
	return &model.Model{
		Logs: log.New(ioutil.Discard, "", 0),
		Controller: &model.Controller{
			Running: false,
		},
	}
}

func MockPauseSend(data string) (*model.Model, string) {
	modl := MockModel()

	res := httptest.NewRecorder()
	req := MockRequest(data)

	Pause(nil, modl)(res, req)

	out, _ := ioutil.ReadAll(res.Body)
	return modl, strings.TrimSpace(string(out))
}

func TestPauseInvalidJSON(t *testing.T) {
	_, res := MockPauseSend("bad json")

	if res != "invalid JSON" {
		t.Fatalf(`expected "invalid JSON", got "%s"`, res)
	}
}

func TestPauseTrue(t *testing.T) {
	modl, res := MockPauseSend(`{"pause":true}`)

	if res != "" {
		t.Fatalf(`expected empty response, got "%s"`, res)
	}

	if !modl.Controller.Paused() {
		t.Fatal("model should be paused")
	}
}

func TestPauseFalse(t *testing.T) {
	modl, res := MockPauseSend(`{"pause":false}`)

	if res != "" {
		t.Fatalf(`expected empty response, got "%s"`, res)
	}

	if modl.Controller.Paused() {
		t.Fatal("model shouldn't be paused")
	}
}

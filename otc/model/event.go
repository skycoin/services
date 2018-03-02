package model

import (
	"time"

	"github.com/skycoin/services/otc/types"
)

type Event struct {
	Time    int64          `json:"time"`
	Request *types.Request `json:"request"`
	Error   string         `json:"error"`
}

func NewEvent(r *types.Request, err error) *Event {
	e := &Event{
		Time:    time.Now().Unix(),
		Request: r,
	}
	if err != nil {
		e.Error = err.Error()
	}
	return e
}

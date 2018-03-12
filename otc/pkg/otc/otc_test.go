package otc

import (
	"testing"
)

func TestRequestId(t *testing.T) {
	req := &Request{
		Address: "1",
		Drop: &Drop{
			Currency: Currency("2"),
			Address:  "3",
		},
	}

	if req.Id() != "1:2:3" {
		t.Fatalf(`expected "1:2:3", got "%s"`, req.Id())
	}
}

func TestRequestIden(t *testing.T) {
	req := &Request{
		Address: "1",
		Drop: &Drop{
			Currency: Currency("2"),
			Address:  "3",
		},
	}

	if req.Iden() != "2:3" {
		t.Fatalf(`expected "2:3", got "%s"`, req.Iden())
	}
}

func TestWorkReturn(t *testing.T) {
	work := &Work{
		Done: make(chan *Result, 1),
	}
	work.Return(nil)

	select {
	case res := <-work.Done:
		if res == nil {
			t.Fatal("nil result")
		}
		return
	default:
		t.Fatal("didn't return")
	}
}

package otc

import (
	"testing"
)

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

func TestOutputsUpdate(t *testing.T) {
	outputs := Outputs(map[string]map[int]*OutputVerbose{
		"transaction": {
			1: {
				Amount: 1,
			},
		},
	})

	outputs.Update("transaction", 1, &OutputVerbose{
		Amount: 2,
	})

	outputs.Update("transaction_two", 1, &OutputVerbose{
		Amount: 3,
	})

	if outputs["transaction"][1].Amount != 2 {
		t.Fatal("update failed")
	}

	if outputs["transaction_two"][1].Amount != 3 {
		t.Fatal("update failed")
	}
}

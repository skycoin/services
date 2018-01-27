package op2c

import (
	"sync"

	"runtime"

	"github.com/skycoin/viscript/signal/msg"
)

type Top struct {
}

type TopResp struct {
	msg.AbstractBlockResp
	runtime.MemStats
}

func init() {
	OPS[OP_TOP] = &sync.Pool{
		New: func() interface{} {
			return new(Top)
		},
	}
	RESPS[OP_TOP] = &sync.Pool{
		New: func() interface{} {
			return new(TopResp)
		},
	}
}

func (r *Top) Execute(c msg.OPer) (resp msg.Resp, err error) {
	result := &TopResp{}
	runtime.ReadMemStats(&result.MemStats)
	resp = result
	return
}

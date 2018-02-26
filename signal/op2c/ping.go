package op2c

import (
	"sync"

	"github.com/skycoin/services/signal/msg"
)

type Ping struct {
}

type PingResp struct {
	msg.AbstractBlockResp
}

func init() {
	OPS[OP_PING] = &sync.Pool{
		New: func() interface{} {
			return new(Ping)
		},
	}
	RESPS[OP_PING] = &sync.Pool{
		New: func() interface{} {
			return new(PingResp)
		},
	}
}

func (r *Ping) Execute(c msg.OPer) (resp msg.Resp, err error) {
	return
}

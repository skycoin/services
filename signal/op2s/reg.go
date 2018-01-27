package op2s

import (
	"sync"

	"github.com/skycoin/viscript/signal/msg"
)

type Reg struct {
	Id uint
}

func init() {
	OPS[OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(Reg)
		},
	}
	RESPS[OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(Reg)
		},
	}
}

func (r *Reg) Execute(c msg.OPer) (resp msg.Resp, err error) {
	c.SetReg(r)
	resp = r
	return
}

func (r *Reg) Receive(c msg.OPer) (err error) {
	c.SetReg(r)
	return
}

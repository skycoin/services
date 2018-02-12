package op2c

import (
	"sync"

	"os"
	"time"

	"github.com/skycoin/services/signal/msg"
)

type Shutdown struct {
}

type ShutdownResp struct {
	msg.AbstractBlockResp
	Pid int
}

func init() {
	OPS[OP_SHUTDOWN] = &sync.Pool{
		New: func() interface{} {
			return new(Shutdown)
		},
	}
	RESPS[OP_SHUTDOWN] = &sync.Pool{
		New: func() interface{} {
			return new(ShutdownResp)
		},
	}
}

func (r *Shutdown) Execute(c msg.OPer) (resp msg.Resp, err error) {
	go func() {
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}()
	resp = &ShutdownResp{Pid: os.Getpid()}
	return
}

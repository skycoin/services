package op2c

import (
	"sync"
)

var (
	OPS   = make([]*sync.Pool, OP_SIZE)
	RESPS = make([]*sync.Pool, OP_SIZE)
)

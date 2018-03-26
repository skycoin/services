package generator

import (
	"sync"

	"github.com/skycoin/services/otc/pkg/otc"
)

type Generator struct {
	sync.Mutex

	Users []*otc.User
}

func New() (*Generator, error) {
	return &Generator{
		Users: make([]*otc.User, 0),
	}, nil
}

func (g *Generator) Add(user *otc.User) *otc.Work {
	return nil
}

package signal

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/viscript/signal/msg"
)

func TestConnect(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	s := NewServer()
	err := s.Listen("localhost:10000")
	if err != nil {
		t.Fatal(err)
	}

	c, err := Connect("localhost:10000", 1)
	if err != nil {
		t.Fatal(err)
	}

	if c.GetReg().Id != 1 {
		t.Fatal("c.GetReg().Id != 1")
	}
	t.Logf("%#v", c.GetReg())
}

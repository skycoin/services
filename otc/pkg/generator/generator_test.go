package generator

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

func TestNew(t *testing.T) {
	gen := New(nil, nil, nil)

	if gen.Users == nil {
		t.Fatal("users map not initialized")
	}
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer

	gen := New(log.New(&buf, "", 0), nil, nil)
	gen.Log("test")

	if buf.String() != "test\n" {
		t.Fatalf("expected %s but got '%s'\n", "test", buf.String())
	}
}

func TestCount(t *testing.T) {
	gen := New(nil, nil, nil)
	gen.Add(&otc.User{})
	gen.Add(&otc.User{})

	if gen.Count() != 2 {
		t.Fatal("wrong count")
	}
}

func TestTick(t *testing.T) {
	task := func(c chan struct{}) Task {
		return func(u *otc.User) (*otc.Order, error) {
			c <- struct{}{}
			return nil, nil
		}
	}

	c := make(chan struct{}, 1)

	gen := New(nil, task(c), nil)
	gen.Add(&otc.User{Id: "1"})
	gen.Tick()

	select {
	case <-c:
		break
	default:
		t.Fatal("didn't run task")
	}
}

func TestDelete(t *testing.T) {
	gen := New(nil, nil, nil)
	user := &otc.User{}

	gen.Add(user)
	gen.Delete(user)

	if gen.Count() != 0 {
		t.Fatal("delete didn't work")
	}
}

func TestTaskGood(t *testing.T) {
	task := func(u *otc.User) (*otc.Order, error) {
		return &otc.Order{
			User:   u,
			Id:     "orderId",
			Status: otc.SEND,
		}, nil
	}

	work := make(chan *otc.Work, 3)
	gen := New(nil, task, work)
	user := &otc.User{}
	gen.Add(user)

	gen.Tick()
	gen.Tick()
	gen.Tick()

	if len(user.Orders) != 3 {
		t.Fatal("orders not appended to user")
	}

	select {
	case <-work:
		break
	default:
		t.Fatal("work not sent over chan")
	}
}

func TestTaskBad(t *testing.T) {
	task := func(u *otc.User) (*otc.Order, error) {
		return nil, fmt.Errorf("bad!")
	}

	var buf bytes.Buffer
	gen := New(log.New(&buf, "", 0), task, nil)
	gen.Add(&otc.User{})
	gen.Tick()

	if buf.String() != "bad!\n" {
		t.Fatal("didn't log error")
	}
}

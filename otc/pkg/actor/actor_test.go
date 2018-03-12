package actor

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

func GoodTask(c chan struct{}) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		c <- struct{}{}
		return true, nil
	}
}

func BadTask() func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		return true, fmt.Errorf("test error")
	}
}

func SkipTask() func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		return false, nil
	}
}

func TestNew(t *testing.T) {
	actor := New(log.New(ioutil.Discard, "", log.Ldate), GoodTask(nil))

	if actor.Work == nil {
		t.Fatal("nil work")
	}
	if actor.Task == nil {
		t.Fatal("nil task")
	}
	if actor.Logs == nil {
		t.Fatal("nil logs")
	}
}

func TestAdd(t *testing.T) {
	actor := New(log.New(ioutil.Discard, "", log.Ldate), GoodTask(nil))
	actor.Add(&otc.Work{})

	if actor.Count() != 1 {
		t.Fatalf("actor count should be 1, but is %d", actor.Count())
	}
}

func TestDelete(t *testing.T) {
	work := &otc.Work{}

	actor := New(log.New(ioutil.Discard, "", log.Ldate), GoodTask(nil))
	actor.Add(work)
	actor.Delete(work)

	if actor.Count() != 0 {
		t.Fatalf("actor count should be 0, but is %d", actor.Count())
	}
}

func TestTask(t *testing.T) {
	notif := make(chan struct{}, 1)
	actor := New(log.New(ioutil.Discard, "", log.Ldate), GoodTask(notif))

	actor.Add(&otc.Work{Done: make(chan *otc.Result, 1)})
	actor.Tick()

	select {
	case <-notif:
		return
	default:
		t.Fatal("task not executing")
	}
}

func TestTaskErr(t *testing.T) {
	work := &otc.Work{Done: make(chan *otc.Result, 1)}
	actor := New(log.New(ioutil.Discard, "", log.Ldate), BadTask())

	actor.Add(work)
	actor.Tick()

	select {
	case res := <-work.Done:
		if res.Err == nil {
			t.Fatal("non err response from work")
		}
		return
	default:
		t.Fatal("task not executing")
	}
}

func TestTaskSkip(t *testing.T) {
	actor := New(log.New(ioutil.Discard, "", log.Ldate), SkipTask())
	actor.Add(&otc.Work{Done: make(chan *otc.Result, 1)})
	actor.Tick()

	if actor.Count() != 1 {
		t.Fatal("task didn't skip deletion")
	}
}

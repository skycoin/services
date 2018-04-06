package model

import (
	"sync"
)

type Controller struct {
	sync.RWMutex

	Running  bool
	Stoppers []chan struct{}
}

func NewController(stoppers []chan struct{}) *Controller {
	return &Controller{Stoppers: stoppers}
}

func (c *Controller) Pause() {
	c.Lock()
	defer c.Unlock()
	c.Running = false
}

func (c *Controller) Unpause() {
	c.Lock()
	defer c.Unlock()
	c.Running = true
}

func (c *Controller) Paused() bool {
	c.RLock()
	defer c.RUnlock()
	return !c.Running
}

func (c *Controller) Stop() {
	for _, s := range c.Stoppers {
		s <- struct{}{}
	}
}

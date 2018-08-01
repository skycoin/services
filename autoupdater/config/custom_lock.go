package config

import "sync"

type CustomLock struct {
	locked bool
	sync.RWMutex
}

func (c *CustomLock) IsLock() bool{
	return c.locked
}

func (c *CustomLock) Lock() {
	c.RLock()
	c.locked=true
	c.RUnlock()
}
func (c *CustomLock) Unlock() {
	c.RLock()
	c.locked=false
	c.RUnlock()
}

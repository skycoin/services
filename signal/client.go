package signal

import (
	"sync"

	"encoding/json"
	"errors"

	"fmt"

	"github.com/skycoin/net/factory"
	"github.com/skycoin/services/signal/msg"
	"github.com/skycoin/services/signal/op2c"
	"github.com/skycoin/services/signal/op2s"
)

type Client struct {
	msg.OPManager
	conn        *factory.Connection
	fieldsMutex sync.RWMutex

	reg     *op2s.Reg
	regCond *sync.Cond

	disconnected chan struct{}

	respChans []chan msg.Resp
}

func newClient(conn *factory.Connection, ops, resps []*sync.Pool) (c *Client) {
	c = &Client{conn: conn, OPManager: msg.NewOPManager(ops, resps)}
	c.regCond = sync.NewCond(c.fieldsMutex.RLocker())
	return
}

func Connect(address string, id uint) (c *Client, err error) {
	c = newClient(nil, op2c.OPS, op2s.RESPS)
	c.reg = &op2s.Reg{Id: id}
	err = c.Connect(address)
	return
}

func (c *Client) Connect(address string) (err error) {
	c.fieldsMutex.RLock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.fieldsMutex.RUnlock()

	f := factory.NewTCPFactory()
	conn, err := f.Connect(address)
	if err != nil {
		return
	}
	c.fieldsMutex.Lock()
	c.conn = conn
	c.disconnected = make(chan struct{})
	c.fieldsMutex.Unlock()
	err = c.Send(op2s.OP_REG, c.reg)
	if err != nil {
		return
	}
	go func() {
		c.readLoop()
		close(c.disconnected)
	}()
	return
}

func (c *Client) SetReg(reg interface{}) {
	r, ok := reg.(*op2s.Reg)
	if !ok {
		return
	}
	c.fieldsMutex.Lock()
	c.reg = r
	c.fieldsMutex.Unlock()
	c.regCond.Broadcast()
}

func (c *Client) GetReg() *op2s.Reg {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	if c.reg == nil {
		c.regCond.Wait()
	}
	return c.reg
}

func (c *Client) Send(op byte, v interface{}) (err error) {
	c.fieldsMutex.RLock()
	conn := c.conn
	c.fieldsMutex.RUnlock()
	if conn == nil {
		err = errors.New("conn == nil")
		return
	}
	js, err := json.Marshal(v)
	if err != nil {
		return
	}
	data := make([]byte, len(js)+1)
	data[msg.MSG_OP_BEGIN] = op
	copy(data[msg.MSG_HEADER_END:], js)
	err = conn.Write(data)
	return
}

func (c *Client) WaitUntilDisconnected() {
	c.fieldsMutex.RLock()
	d := c.disconnected
	c.fieldsMutex.RUnlock()
	if d == nil {
		return
	}
	<-d
}

func (c *Client) readLoop() {
	var err error
	conn := c.conn
	defer func() {
		if e := recover(); e != nil {
			conn.GetContextLogger().Errorf("readLoop recover err %v", e)
		}
		if err != nil {
			conn.GetContextLogger().Errorf("readLoop err %v", err)
		}
		c.Close()
	}()
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			err = c.Operate(c, m)
			if err != nil {
				conn.GetContextLogger().Errorf("execute err: %v", err)
			}
		}
	}
}

func (c *Client) Close() {
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	if c.respChans != nil {
		for _, c := range c.respChans {
			close(c)
		}
		c.respChans = nil
	}

	c.regCond.Broadcast()
}

// Server side operations

func (c *Client) Ping() (resp *op2c.PingResp, err error) {
	err = c.Send(op2c.OP_PING, nil)
	if err != nil {
		return
	}
	c.fieldsMutex.RLock()
	ch := c.respChans[op2c.OP_PING]
	c.fieldsMutex.RUnlock()
	resp = (<-ch).(*op2c.PingResp)
	return
}

func (c *Client) Top() (resp *op2c.TopResp, err error) {
	err = c.Send(op2c.OP_TOP, nil)
	if err != nil {
		return
	}
	c.fieldsMutex.RLock()
	ch := c.respChans[op2c.OP_TOP]
	c.fieldsMutex.RUnlock()
	resp = (<-ch).(*op2c.TopResp)
	return
}

func (c *Client) Shutdown() (resp *op2c.ShutdownResp, err error) {
	err = c.Send(op2c.OP_SHUTDOWN, nil)
	if err != nil {
		return
	}
	c.fieldsMutex.RLock()
	ch := c.respChans[op2c.OP_SHUTDOWN]
	c.fieldsMutex.RUnlock()
	resp = (<-ch).(*op2c.ShutdownResp)
	return
}

func (c *Client) initRespChan() {
	respChans := make([]chan msg.Resp, len(c.RespPools))
	for i := range c.RespPools {
		respChans[i] = make(chan msg.Resp)
	}
	c.fieldsMutex.Lock()
	c.respChans = respChans
	c.fieldsMutex.Unlock()
}

func (c *Client) ReceiveBlockResp(op int, resp msg.Resp) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("ReceiveBlockResp recover %v", e)
		}
	}()
	c.fieldsMutex.RLock()
	ch := c.respChans[op]
	c.fieldsMutex.RUnlock()
	ch <- resp
	return
}

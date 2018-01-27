package msg

type OP interface {
	Execute(OPer) (Resp, error)
}

type Resp interface {
	Receive(OPer) error
}

type OPer interface {
	SetReg(interface{})
	Send(byte, interface{}) error
	ReceiveBlockResp(int, Resp) error
}

type BlockResp interface {
	Block()
}

type AbstractBlockResp struct {
}

func (r *AbstractBlockResp) Receive(c OPer) (err error) {
	return
}

func (r *AbstractBlockResp) Block() {
}

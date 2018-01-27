package msg

import (
	"encoding/json"
	"sync"

	"fmt"

	log "github.com/sirupsen/logrus"
)

type OPManager struct {
	opPools   []*sync.Pool
	RespPools []*sync.Pool
}

func NewOPManager(op, resp []*sync.Pool) (opm OPManager) {
	return OPManager{opPools: op, RespPools: resp}
}

func (opm *OPManager) getOP(opn byte) interface{} {
	if opn&RESP_PREFIX > 0 {
		opn = opn &^ RESP_PREFIX
		if opn < 0 || int(opn) > len(opm.RespPools) {
			return nil
		}

		return opm.RespPools[opn].Get()
	}

	if opn < 0 || int(opn) > len(opm.opPools) {
		return nil
	}

	return opm.opPools[opn].Get()
}

func (opm *OPManager) putOP(opn byte, op interface{}) {
	if opn&RESP_PREFIX > 0 {
		opn = opn &^ RESP_PREFIX
		if opn < 0 || int(opn) > len(opm.RespPools) {
			return
		}
		opm.RespPools[opn].Put(op)
		return
	}
	if opn < 0 || int(opn) > len(opm.opPools) {
		return
	}
	opm.opPools[opn].Put(op)
}

func (opm *OPManager) Operate(oper OPer, m []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("operate err %v", e)
		}
	}()
	if len(m) < MSG_HEADER_END {
		return
	}
	opn := m[MSG_OP_BEGIN]
	op := opm.getOP(opn)
	if op == nil {
		return
	}
	defer func() {
		opm.putOP(opn, op)
	}()

	err = json.Unmarshal(m[MSG_HEADER_END:], op)
	if err == nil {
		if opn&RESP_PREFIX > 0 {
			resp, ok := op.(Resp)
			if !ok {
				return fmt.Errorf("%#v can not convert to Resp", op)
			}
			log.Debugf("receive resp %#v", op)
			if _, isBlock := resp.(BlockResp); isBlock {
				err = oper.ReceiveBlockResp(int(opn&^RESP_PREFIX), resp)
			} else {
				err = resp.Receive(oper)
			}
		} else {
			op, ok := op.(OP)
			if !ok {
				return fmt.Errorf("%#v can not convert to OP", op)
			}
			log.Debugf("receive op %#v", op)
			var resp Resp
			resp, err = op.Execute(oper)
			if err != nil {
				return
			}
			err = oper.Send(opn|RESP_PREFIX, resp)
		}
	}
	return
}

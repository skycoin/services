package sender

import (
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(curs *currencies.Currencies) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		work.Request.Lock()
		defer work.Request.Unlock()

		value, source, err := curs.Value(work.Request.Drop)
		if err != nil {
			return true, err
		}
		work.Request.Rate = &otc.Rate{Value: value, Source: source}

		txid, err := curs.Send(otc.SKY, work.Request.Address, value)
		if err != nil {
			return true, err
		}

		work.Request.Times.SentAt = time.Now().UTC().Unix()
		work.Request.TxId = txid
		work.Request.Status = otc.CONFIRM
		return true, nil
	}
}

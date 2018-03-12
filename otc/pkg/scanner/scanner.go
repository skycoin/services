package scanner

import (
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(curs *currencies.Currencies) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		work.Request.Lock()
		defer work.Request.Unlock()

		balance, err := curs.Balance(work.Request.Drop)
		if err != nil {
			return true, err
		}

		if balance != 0 {
			work.Request.Times.DepositedAt = time.Now().UTC().Unix()
			work.Request.Drop.Amount = balance
			work.Request.Status = otc.SEND
			return true, nil
		}

		return false, nil
	}
}

package monitor

import (
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(curs *currencies.Currencies) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		work.Request.Lock()
		defer work.Request.Unlock()

		confirmed, err := curs.Confirmed(otc.SKY, work.Request.TxId)
		if err != nil {
			return true, err
		}

		if confirmed {
			work.Request.Times.ConfirmedAt = time.Now().UTC().Unix()
			work.Request.Status = otc.DONE
			return true, nil
		}

		return false, nil
	}
}

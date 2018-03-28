package monitor

import (
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(curs *currencies.Currencies) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		confirmed, err := curs.Confirmed(otc.SKY, work.Order.Purchase.TxId)
		if err != nil {
			return true, err
		}

		if confirmed {
			work.Order.Times.ConfirmedAt = time.Now().UTC().Unix()
			work.Order.Status = otc.DONE
			return true, nil
		}

		return false, nil
	}
}

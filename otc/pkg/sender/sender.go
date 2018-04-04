package sender

import (
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(curs *currencies.Currencies) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		value, source, price, err := curs.Value(
			work.Order.User.Drop.Currency,
			work.Order.Amount,
		)
		if err != nil {
			return true, err
		}

		txid, err := curs.Send(otc.SKY, work.Order.User.Address, value)
		if err != nil {
			return true, err
		}

		work.Order.Purchase = &otc.Purchase{
			// TODO: make source string dynamic
			Source: "internal",
			Amount: value,
			TxId:   txid,
			Price:  &otc.Price{source, price},
		}
		work.Order.Times.SentAt = time.Now().UTC().Unix()
		work.Order.Status = otc.CONFIRM
		return true, nil
	}
}

package scanner

import (
	"fmt"
	"time"

	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/services/otc/pkg/watcher"
)

func Task(watch *watcher.Watcher) func(*otc.User) (*otc.Order, error) {
	return func(user *otc.User) (*otc.Order, error) {
		// get deposits from otc-watcher
		deposits, err := watch.Outputs(user.Drop)
		if err != nil {
			return nil, err
		}

		for transaction, outputs := range deposits {
		indexing:
			for index, output := range outputs {
				id := fmt.Sprintf("%s:%d", transaction, index)

				// check if order already exists
				for _, order := range user.Orders {
					if order.Id == id {
						continue indexing
					}
				}

				now := time.Now().UTC().Unix()

				// generate new order
				return &otc.Order{
					Id:     id,
					Status: otc.SEND,
					Amount: output.Amount,
					Times: &otc.Times{
						CreatedAt:   now,
						UpdatedAt:   now,
						DepositedAt: now,
					},
					Events: []*otc.Event{
						{
							Status:   otc.DEPOSIT,
							Finished: now,
						},
					},
				}, nil
			}
		}

		return nil, nil
	}
}

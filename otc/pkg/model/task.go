package model

import (
	"time"

	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(workers *Workers) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		select {
		case res := <-work.Done:
			work.Order.Times.UpdatedAt = time.Now().UTC().Unix()

			// save to disk
			if err := SaveOrder(work.Order, res); err != nil {
				return true, err
			}

			// check result
			if res.Err != nil {
				return true, res.Err
			}

			// if done, stop routing
			if work.Order.Status == otc.DONE {
				return true, nil
			}

			// route to next step
			workers.Route(work)
		default:
			break
		}

		return false, nil
	}
}

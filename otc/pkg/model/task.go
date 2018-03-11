package model

import (
	"time"

	"github.com/skycoin/services/otc/pkg/otc"
)

func Task(workers *Workers) func(*otc.Work) (bool, error) {
	return func(work *otc.Work) (bool, error) {
		work.Request.Lock()
		defer work.Request.Unlock()

		select {
		case res := <-work.Done:
			work.Request.Times.UpdatedAt = time.Now().UTC().Unix()

			if err := Save(work.Request, res); err != nil {
				return true, err
			}

			if res.Err != nil {
				return true, res.Err
			}

			if work.Request.Status == otc.DONE {
				return true, nil
			}

			workers.Route(work)
		default:
			break
		}

		return false, nil
	}
}

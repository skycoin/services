package scanner

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/skycoin/services/otc-watcher/pkg/currency"
	"github.com/skycoin/services/otc/pkg/otc"
)

func (s *Scanner) Save(cur otc.Currency) error {
	file, err := os.OpenFile(string(cur)+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	s.Scanning[cur].RLock()
	defer s.Scanning[cur].RUnlock()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err = enc.Encode(s.Scanning[cur]); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return file.Close()
}

func (s *Scanner) Load(cons currency.Connections) error {
	for cur, _ := range cons {
		file, err := os.OpenFile(string(cur)+".json", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		s.Scanning[cur] = NewStorage(cur)

		if err = json.NewDecoder(file).Decode(s.Scanning[cur]); err != nil {
			if err == io.EOF {
				// get latest block height
				height, err := cons.Height(cur)
				if err != nil {
					return err
				}

				// create and start scanning from latest block
				s.Scanning[cur].Updated = &Updated{
					Time:   time.Now().UTC().Unix(),
					Height: height,
				}
				s.Scanning[cur].Addresses = make(map[string]*Relevant, 0)
			} else {
				return err
			}
		}

		if err = file.Close(); err != nil {
			return err
		}
	}

	return nil
}

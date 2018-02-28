package model

import (
	"errors"
	"os"
	"sync"

	"github.com/skycoin/services/otc/types"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	STORAGE_REQUESTS = "requests/"
)

type Storage struct {
	sync.RWMutex

	Path string
}

func NewStorage(path string) (*Storage, error) {
	s := &Storage{
		Path: path,
	}

	// check that storage path exists
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return s, nil
}

var ErrInvalidFilename = errors.New("invalid filename in db dir")

func (s *Storage) LoadRequests(path string) ([]*types.Request, error) {
	// check that filename is longer than just ".json"
	if len(path) <= 5 {
		return nil, ErrInvalidFilename
	}

	// check that filename is a valid sky address
	addr, err := cipher.DecodeBase58Address(path[:len(path)-5])
	if err != nil {
		return nil, err
	}

	// read lock for this address
	s.RLock()
	defer s.RUnlock()

	// get raw map from disk
	data, err := mapFromJSON(s.Path + STORAGE_REQUESTS + path)
	if err != nil {
		return nil, err
	}

	requests := make([]*types.Request, 0)

	for currency, drops := range data {
		for drop, metadata := range drops {
			// ignore finished requests
			if metadata.Status == types.DONE {
				continue
			}
			requests = append(requests, &types.Request{
				Address:  types.Address(addr.String()),
				Currency: types.Currency(currency),
				Drop:     types.Drop(drop),
				Metadata: metadata,
			})
		}
	}

	return requests, nil
}

var ErrDropMissing = errors.New("drop doesn't exist")

func (s *Storage) LoadMetadata(address types.Address, drop types.Drop, curr types.Currency) (*types.Metadata, error) {
	s.RLock()
	defer s.RUnlock()

	// full filepath for .json file
	path := s.Path + STORAGE_REQUESTS + string(address) + ".json"

	// read json from disk
	data, err := mapFromJSON(path)
	if err != nil {
		return nil, err
	}

	// check that the currency type and drop exists
	if data[curr] == nil || data[curr][drop] == nil {
		return nil, ErrDropMissing
	}

	// return metadata from file
	return data[curr][drop], nil
}

func (s *Storage) SaveRequest(request *types.Request) error {
	s.Lock()
	defer s.Unlock()

	// full filepath for .json file
	path := s.Path + STORAGE_REQUESTS + string(request.Address) + ".json"

	// read json data from disk
	data, err := mapFromJSON(path)
	if err != nil {
		return err
	}

	// update map
	if data == nil {
		data = map[types.Currency]map[types.Drop]*types.Metadata{
			request.Currency: {request.Drop: request.Metadata},
		}
	} else if data[request.Currency] == nil {
		data[request.Currency] = map[types.Drop]*types.Metadata{
			request.Drop: request.Metadata,
		}
	} else {
		data[request.Currency][request.Drop] = request.Metadata
	}

	// write map to disk
	return mapToJSON(path, data)
}

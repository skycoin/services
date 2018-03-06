package model

import (
	"encoding/json"
	"io"
	"os"

	"github.com/skycoin/services/otc/types"
)

func mapFromJSON(path string) (map[types.Currency]map[types.Drop]*types.Metadata, error) {
	empty := make(map[types.Currency]map[types.Drop]*types.Metadata)

	// read json file from disk
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if err == io.EOF {
			return empty, nil
		}
		return nil, err
	}

	// will be returned
	var data map[types.Currency]map[types.Drop]*types.Metadata

	// decode json from file
	if err = json.NewDecoder(file).Decode(&data); err != nil {
		if err == io.EOF {
			return empty, nil
		}
		return nil, err
	}

	return data, file.Close()
}

func mapToJSON(path string, data map[types.Currency]map[types.Drop]*types.Metadata) error {
	// open file for writing
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	// reset file
	file.Truncate(0)
	file.Seek(0, 0)

	// json encoder with indent
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// encode to file
	if err = enc.Encode(data); err != nil {
		return err
	}

	// sync file contents to disk
	if err = file.Sync(); err != nil {
		return err
	}

	return file.Close()
}

package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	PATH string = ".otc/"
	REQS string = "reqs/"
	LOGS string = "logs/"
)

func Save(req *otc.Request, res *otc.Result) error {
	file, err := os.OpenFile(
		PATH+REQS+req.Id()+".json",
		os.O_CREATE|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}

	// empty file
	file.Truncate(0)
	file.Seek(0, 0)

	// indent json
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// write json to file
	if err = enc.Encode(req); err != nil {
		return err
	}

	// sync to disk
	if err = file.Sync(); err != nil {
		return err
	}

	// close file
	if err = file.Close(); err != nil {
		return err
	}

	return Log(req, res)
}

func Log(req *otc.Request, res *otc.Result) error {
	file, err := os.OpenFile(
		PATH+LOGS+"log.json",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	event := &otc.Event{
		Id:       req.Id(),
		Status:   req.Status,
		Finished: res.Finished,
	}
	if res.Err != nil {
		event.Err = res.Err.Error()
	}

	if err = json.NewEncoder(file).Encode(&event); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return file.Close()
}

func Load() ([]*otc.Request, error) {
	// get list of files in db dir
	files, err := ioutil.ReadDir(PATH + REQS)
	if err != nil {
		return nil, err
	}

	reqs := make([]*otc.Request, 0)

	// for each .json file in db dir
	for _, file := range files {
		// ignore hidden files
		if file.Name()[0] == '.' {
			continue
		}

		// get struct from json
		req, err := Read(PATH+REQS, file.Name())
		if err != nil {
			return nil, err
		}

		// append to slice
		reqs = append(reqs, req)
	}

	return reqs, nil
}

func Read(path, filename string) (*otc.Request, error) {
	parts := strings.Split(filename, ":")

	// check that filename is in form of x:x:x
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request filename")
	}

	// check that first part is valid sky address
	_, err := cipher.DecodeBase58Address(parts[0])
	if err != nil {
		return nil, err
	}

	// open file for reading
	file, err := os.OpenFile(path+filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var req *otc.Request

	if err = json.NewDecoder(file).Decode(&req); err != nil {
		return nil, err
	}

	return req, file.Close()
}

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
	PATH   string = ".otc/"
	USERS  string = "users/"
	ORDERS string = "orders/"
)

func SaveUser(user *otc.User) error {
	file, err := os.OpenFile(
		PATH+USERS+user.Id+".json",
		os.O_CREATE|os.O_RDWR, 0644,
	)
	if err != nil {
		return err
	}

	// empty file if exists
	file.Truncate(0)
	file.Seek(0, 0)

	// format json output
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err = enc.Encode(user); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	// create orders folder
	return os.Mkdir(PATH+ORDERS+user.Id, os.ModeDir)
}

func SaveOrder(order *otc.Order, result *otc.Result) error {
	file, err := os.OpenFile(
		PATH+ORDERS+order.User.Id+"/"+order.Id+".json",
		os.O_CREATE|os.O_RDWR, 0644,
	)
	if err != nil {
		return err
	}

	// append to order events
	event := &otc.Event{
		Status:   order.Status,
		Finished: result.Finished,
	}
	if result.Err != nil {
		event.Err = result.Err.Error()
	}
	order.Events = append(order.Events, event)

	// empty file
	file.Truncate(0)
	file.Seek(0, 0)

	// indent json output
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err = enc.Encode(order); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	return file.Close()
}

func Load() ([]*otc.User, error) {
	// get list of users
	files, err := ioutil.ReadDir(PATH + USERS)
	if err != nil {
		return nil, err
	}

	// for returning later
	users := make([]*otc.User, 0)

	for _, file := range files {
		// ignore hidden files
		if file.Name()[0] == '.' {
			continue
		}

		// get user struct from disk
		user, err := ReadUser(PATH+USERS, file.Name())
		if err != nil {
			return nil, err
		}

		// get list of orders in user's dir
		ofiles, err := ioutil.ReadDir(PATH + ORDERS + user.Id)
		if err != nil {
			return nil, err
		}

		// for each order file associated with user
		for _, ofile := range ofiles {
			// ignore hidden files
			if ofile.Name()[0] == '.' {
				continue
			}

			// read order from disk
			order, err := ReadOrder(PATH+ORDERS+user.Id, ofile.Name())
			if err != nil {
				return nil, err
			}

			// add order to user
			user.Orders = append(user.Orders, order)
		}

		// append to list of all users
		users = append(users, user)
	}

	return users, nil
}

func ReadUser(path, filename string) (*otc.User, error) {
	parts := strings.Split(filename, ":")

	// check that filename is in form of x:x:x
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid user filename")
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

	user := new(otc.User)

	// read from disk
	if err = json.NewDecoder(file).Decode(&user); err != nil {
		return nil, err
	}

	return user, file.Close()
}

func ReadOrder(path, filename string) (*otc.Order, error) {
	parts := strings.Split(filename, ":")

	// check that filename is in form of x:x
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid order filename")
	}

	// open file for reading
	file, err := os.OpenFile(path+filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	order := new(otc.Order)

	// read from disk
	if err = json.NewDecoder(file).Decode(&order); err != nil {
		return nil, err
	}

	return order, file.Close()
}

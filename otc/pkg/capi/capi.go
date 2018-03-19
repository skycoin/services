package capi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/services/otc/pkg/otc"
)

type Client struct {
	Base   string
	Client *http.Client
}

func New(conf *otc.Config) *Client {
	return &Client{
		Base: conf.CAPI.Node,
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *Client) Unspents(curr otc.Currency, addr string) (otc.Unspents, error) {
	var buf bytes.Buffer

	// encode request into json
	err := json.NewEncoder(&buf).Encode(struct {
		Address string `json:"address"`
	}{addr})
	if err != nil {
		return nil, err
	}

	// send post request to coin-api
	res, err := c.Client.Post(c.Base+"/unspents", "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// create output struct to deserialize json into
	//
	// TODO: make this a global type w/ wrapper functions
	out := &struct {
		Status  string       `json:"status"`
		Code    int          `json:"code"`
		Message string       `json:"message"`
		Result  otc.Unspents `json:"result"`
	}{}

	// attempt to decode json
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	// check that coin-api didn't return an error
	if out.Status != "Ok" {
		return nil, fmt.Errorf(out.Message)
	}

	return out.Result, nil
}

func (c *Client) Register(curr otc.Currency, addr string) error {
	var buf bytes.Buffer

	// encode request into json
	err := json.NewEncoder(&buf).Encode(struct {
		Address string `json:"address"`
	}{addr})
	if err != nil {
		return err
	}

	// send post request to coin-api
	res, err := c.Client.Post(c.Base+"/register", "application/json", &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// create output struct to deserialize json into
	//
	// TODO: make this a global type w/ wrapper functions
	out := &struct {
		Status  string `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}

	// attempt to decode json
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return err
	}

	// check that coin-api didn't return an error
	if out.Status != "Ok" {
		return fmt.Errorf(out.Message)
	}

	return nil
}

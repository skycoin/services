package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/skycoin/services/bitcoin-scanning-wallet/config"
	"github.com/skycoin/services/bitcoin-scanning-wallet/scan"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func responseMessage(data interface{}, response http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(js)
}

//ServerStatus is struct
type ServerStatus struct {
	Alive bool
	Time  time.Time
}

//StatusHandler is used for check server state
var StatusHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	var status ServerStatus
	status.Alive = true
	status.Time = time.Now()
	responseMessage(status, response, request)
})

//AddressHandler show all wallet
var AddressHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}
	responseMessage(addrs, response, request)
})

//MinScanHandler scan min block
var MinScanHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}

	client, err := scan.NewBTCDClient(config.Config.BTCD.User, config.Config.BTCD.Pass)
	defer client.Shutdown()

	addrs, err = scan.UpdateMin(addrs, client)
	if err != nil {
		fmt.Println("Update min is failed:", err)
	}

	err = scan.SaveWallet(config.Config.Wallet.File, addrs)
	if err != nil {
		fmt.Println("Saving wallet is failed:", err)
	}
	responseMessage(addrs, response, request)
})

//MaxScanHandler scan min block
var MaxScanHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}

	client, err := scan.NewBTCDClient(config.Config.BTCD.User, config.Config.BTCD.Pass)
	defer client.Shutdown()

	addrs, err = scan.UpdateMax(addrs, client)
	if err != nil {
		fmt.Println("Update max is failed:", err)
	}

	err = scan.SaveWallet(config.Config.Wallet.File, addrs)
	if err != nil {
		fmt.Println("Saving wallet is failed:", err)
	}
	responseMessage(addrs, response, request)
})

//FarScanHandler scan fat block
var FarScanHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}

	client, err := scan.NewBTCDClient(config.Config.BTCD.User, config.Config.BTCD.Pass)
	defer client.Shutdown()

	addrs, err = scan.UpdateFar(addrs, client)
	if err != nil {
		fmt.Println("Update far is failed:", err)
	}

	err = scan.SaveWallet(config.Config.Wallet.File, addrs)
	if err != nil {
		fmt.Println("Saving wallet is failed:", err)
	}
	responseMessage(addrs, response, request)
})

//ShortScanHandler scan short block
var ShortScanHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}

	client, err := scan.NewBTCDClient(config.Config.BTCD.User, config.Config.BTCD.Pass)
	defer client.Shutdown()

	addrs, err = scan.UpdateShort(addrs, client)
	if err != nil {
		fmt.Println("Update short is failed:", err)
	}

	err = scan.SaveWallet(config.Config.Wallet.File, addrs)
	if err != nil {
		fmt.Println("Saving wallet is failed:", err)
	}
	responseMessage(addrs, response, request)
})

//DiapasonHandler get start and end block number and scan these blocks
var DiapasonHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {

	}
	n, _ := strconv.Atoi(request.FormValue("n"))
	m, _ := strconv.Atoi(request.FormValue("m"))
	fmt.Println(n, m)

	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}

	//create btcd instance
	client, err := scan.NewBTCDClient(config.Config.BTCD.User, config.Config.BTCD.Pass)
	defer client.Shutdown()

	for i := n; i <= m; i++ {
		//fmt.Println("Scannig block: ", i)
		deposits, err := scan.ScanBlock(client, int64(i))
		if err != nil {
			fmt.Println("Block scanning is failed:", err)
		}

		addrs = scan.UpdateAddressInfo(addrs, deposits, int64(i))
	}

	err = scan.SaveWallet(config.Config.Wallet.File, addrs)
	if err != nil {
		fmt.Println("Saving wallet is failed:", err)
	}

	responseMessage(addrs, response, request)
})

//AddAddressHandler add address to wallet
var AddAddressHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	json, _ := gabs.ParseJSON(body)
	children, _ := json.S("addrs").Children()
	for _, child := range children {
		fmt.Println(child.Data().(string))
		scan.AddBTCAddress(child.Data().(string), config.Config.Wallet.File)
	}
	responseMessage(json, response, request)

})

//AddAddressHandler get address and return all transactions
var GetAddressHandler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	address := request.URL.Query().Get("address")
	addrs, err := scan.LoadWallet(config.Config.Wallet.File)
	if err != nil {
		fmt.Println("Wallet loading is failed:", err)
	}
	var newAddr scan.Address
	_, newAddr = scan.FindAddress(address, addrs)

	responseMessage(newAddr, response, request)

})

package main

import (
	"net/http"
	"./handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"./scan"
	"./config"
	"fmt"
	"os"
)


func main() {

	config := config.LoadConfiguration("config.json")
	_, err := scan.NewBTCDClient(config.BTCD.User, config.BTCD.Pass)
	if err != nil {
		fmt.Printf("Can't connect btcd, error: ", err)
		os.Exit(1)
	} else {
		fmt.Println("Connect to btcd is established")
	}
	startServer()

}

func startServer() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.Handle("/getaddrs", handler.AddressHandler).Methods("GET")
	r.Handle("/scanrange", handler.DiapasonHandler).Methods("POST")
	r.Handle("/newaddrs", handler.AddAddressHandler).Methods("POST")
	http.ListenAndServe(":7755", handlers.CORS()(r))
}



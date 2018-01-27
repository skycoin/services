package signal

import (
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// signal server address
	serverAddress string = "localhost:7999"
	// run as a signal client
	runClient bool
	// client id
	clientId uint = 1
)

func init() {
	// look up env values
	value, ok := os.LookupEnv("SIGNAL_CLIENT")
	if ok {
		if strings.ToLower(value) != "false" {
			runClient = true
		}
	}

	value, ok = os.LookupEnv("SIGNAL_CLIENT_ID")
	if ok {
		id, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			clientId = uint(id)
		}
	}

	value, ok = os.LookupEnv("SIGNAL_SERVER_ADDRESS")
	if ok {
		serverAddress = value
	}

	if runClient {
		go func() {
			c, err := Connect(serverAddress, clientId)
			for {
				if err != nil {
					log.Errorf("connect to viscript as id %d failed %v", clientId, err)
				}
				c.WaitUntilDisconnected()
				// sleep 30s to reconnect
				time.Sleep(30 * time.Second)
				err = c.Connect(serverAddress)
			}
		}()
	}
}

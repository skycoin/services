package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

func main() {
	coin := flag.String("c", "SKY", "code of new coin")
	file := flag.String("f", "", "file to save configuration of new coin")

	flag.Parse()

	cfg := createCoin(*coin)

	// Print config
	out, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal JSON - %s", err)
	}

	fmt.Println(string(out))

	// Save config to disk if required
	if *file != "" {
		if err := ioutil.WriteFile(*file, out, os.ModePerm); err != nil {
			log.Fatalf("failed to save coin configuration to file - %s", err)
		}
	}
}

func createCoin(coinCode string) common.Config {
	// Generate master's private and public key pair
	pkb := make([]byte, 32)
	if _, err := rand.Read(pkb); err != nil {
		log.Fatalf("failed to create master's private key")
	}

	sk := cipher.NewSecKey(pkb)
	pk := cipher.PubKeyFromSecKey(sk)

	// Geneate genesis block
	var (
		gbAddr  = cipher.AddressFromSecKey(sk)
		gbCoins = uint64(100e12)
		gbTs    = uint64(time.Now().Unix())
	)
	gb, err := coin.NewGenesisBlock(gbAddr, gbCoins, gbTs)
	if err != nil {
		log.Fatalf("failed to create genesis block - %s", err)
	}

	// Coin configuration
	cfg := common.Config{
		Secret: common.SecretConfig{
			MasterSecKey:     sk.Hex(),
			GenesisSignature: cipher.SignHash(gb.HashHeader(), sk).Hex(),
		},

		Public: common.PublicConfig{
			MasterPubKey: pk.Hex(),

			GenesisBlock: common.GenesisBlockConfig{
				Address:    gbAddr.String(),
				CoinVolume: gbCoins,
				Timestamp:  gbTs,
				BodyHash:   gb.HashBody().Hex(),
				HeaderHash: gb.HashHeader().Hex(),
			},

			CoinCode: coinCode,

			Port:             16000,
			WebInterfacePort: 16420,
			RPCInterfacePort: 16430,

			DataDirectory: ".skycoin-testnet",
			LogFmt:        "[skycoin.testnet.%{module}:%{level}] %{message}",
		},
	}

	return cfg
}

package main

import (
	"encoding/hex"
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
	var (
		file      = flag.String("file", "", "file to save configuration of new coin")
		coin      = flag.String("code", "SKY", "code of new coin")
		addrCount = flag.Int("addr", 100, "number of distribution addresses")
		coinVol   = flag.Int("vol", int(1e6), "coin volume to send to each of disribution addresses")
	)

	flag.Parse()

	cfg := createCoin(*coin, *addrCount, *coinVol)

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

func createCoin(coinCode string, addrCount, coinVol int) common.Config {
	sk := cipher.NewSecKey(cipher.RandByte(32))
	pk := cipher.PubKeyFromSecKey(sk)

	// Geneate genesis block
	var (
		gbAddr  = cipher.AddressFromSecKey(sk)
		gbCoins = uint64(addrCount * coinVol)
		gbTs    = uint64(time.Now().Unix())
	)
	gb, err := coin.NewGenesisBlock(gbAddr, gbCoins, gbTs)
	if err != nil {
		log.Fatalf("failed to create genesis block - %s", err)
	}

	addrSeed, addrs := genDistAdrresses(addrCount)

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

			Distribution: common.DistributionConfig{
				CoinsPerAddress: uint64(coinVol),
				AddressSeed:     addrSeed,
				Addresses:       addrs,
			},

			Port:             16000,
			WebInterfacePort: 16420,
			RPCInterfacePort: 16430,
		},
	}

	return cfg
}

func genDistAdrresses(n int) (string, []string) {
	seed := cipher.RandByte(64)

	addrs := make([]string, n)
	keys := cipher.GenerateDeterministicKeyPairs(seed, n)
	for i, k := range keys {
		addrs[i] = cipher.AddressFromSecKey(k).String()
	}

	return hex.EncodeToString(seed), addrs
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/coin"
)

const (
	trustedPeerPort = 20000
	daemonPort      = 20100
	rpcPort         = 20200
	guiPort         = 20300
)

func main() {
	var (
		file      = flag.String("file", "", "file to save configuration of new coin")
		coin      = flag.String("code", "SKY", "code of new coin")
		addrCount = flag.Int("addr", 100, "number of distribution addresses")
		coinVol   = flag.Int("vol", 1, "coin volume to send to each of disribution addresses")
		peerCount = flag.Int("peers", 3, "number of trusted peers running on localhost")
	)

	flag.Parse()

	cfg := createCoin(*coin, *addrCount, *coinVol, *peerCount)

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

func createCoin(coinCode string, addrCount, coinVol, peerCount int) common.Config {
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

	// Genesis block wallet
	gwSeed, err := bip39.NewDefaultMnemomic()
	if err != nil {
		log.Fatalf("failed to generate genesis wallet seed")
	}

	// Trusted peers of coin networks (default connections)
	peers := make([]string, peerCount)
	for i := 0; i < peerCount; i++ {
		peers[i] = fmt.Sprintf("127.0.0.1:%d", trustedPeerPort+i)
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

			GenesisWallet: common.GenesisWalletConfig{
				Seed:            gwSeed,
				CoinsPerAddress: uint64(coinVol),
				Addresses:       uint64(addrCount),
			},

			CoinCode: coinCode,

			Port:    daemonPort,
			RPCPort: rpcPort,
			GUIPort: guiPort,

			TrustedPeers: peers,
		},
	}

	return cfg
}

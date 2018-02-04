package main

import (
	"errors"
	"time"

	"github.com/skycoin/skycoin/src/coin"

	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
)

// Config records the node's configuration
type NodeConfig struct {
	// Disable peer exchange
	DisablePEX bool
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max connections to maintain
	MaxConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface      bool
	WebInterfacePort  int
	WebInterfaceAddr  string
	WebInterfaceCert  string
	WebInterfaceKey   string
	WebInterfaceHTTPS bool

	RPCInterface     bool
	RPCInterfacePort int
	RPCInterfaceAddr string

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// GUI directory contains assets for the html gui
	GUIDirectory string

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string

	RunMaster bool

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

	DefaultConnections []string

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo string

	DBPath       string
	Arbitrating  bool
	RPCThreadNum uint   // rpc number
	LogFmt       string // log format
	Logtofile    bool
	TestChain    bool
}

func makeNodeConfig(toolCfg common.Config) (NodeConfig, error) {
	var (
		cfg NodeConfig
		err error
	)

	// Hardcoded configuration
	cfg.MaxConnections = 16
	cfg.OutgoingConnectionsRate = time.Second * 5
	cfg.WebInterface = true
	cfg.WebInterfaceAddr = "127.0.0.1"
	cfg.RPCInterface = true
	cfg.RPCInterfaceAddr = "127.0.0.1"
	cfg.RPCThreadNum = 5
	cfg.LaunchBrowser = true
	cfg.GUIDirectory = "./src/gui/static/"
	cfg.ColorLog = true
	cfg.LogLevel = "DEBUG"
	cfg.ProfileCPUFile = "skycoin.prof"

	// User provided configuration
	if cfg.BlockchainSeckey, err = cipher.SecKeyFromHex(toolCfg.Secret.MasterSecKey); err != nil {
		return cfg, errors.New("invalid master node secret key")
	}

	if cfg.BlockchainPubkey, err = cipher.PubKeyFromHex(toolCfg.Public.MasterPubKey); err != nil {
		return cfg, errors.New("invalid master node public key")
	}

	cfg.Port = toolCfg.Public.Port
	cfg.WebInterfacePort = toolCfg.Public.WebInterfacePort
	cfg.RPCInterfacePort = toolCfg.Public.RPCInterfacePort

	cfg.DataDirectory = toolCfg.Public.DataDirectory
	// TODO
	// cfg.DefaultConnections = toolCfg.Public.DefaultConnections
	cfg.LogFmt = toolCfg.Public.LogFmt

	return cfg, nil
}

func makeGenesisBlock(cfg common.Config) (coin.SignedBlock, error) {
	sig, err := cipher.SigFromHex(cfg.Secret.GenesisSignature)
	if err != nil {
		return coin.SignedBlock{}, errors.New("invalid genesis block signature")
	}

	addr, err := cipher.DecodeBase58Address(cfg.Public.GenesisBlock.Address)
	if err != nil {
		return coin.SignedBlock{}, errors.New("invalid genesis block address")
	}

	bodyHash, err := cipher.SHA256FromHex(cfg.Public.GenesisBlock.BodyHash)
	if err != nil {
		return coin.SignedBlock{}, errors.New("invalid genesis block body hash")
	}

	var tx coin.Transaction
	tx.PushOutput(addr, cfg.Public.GenesisBlock.CoinVolume, cfg.Public.GenesisBlock.CoinVolume)

	b := coin.Block{
		Head: coin.BlockHeader{
			Time:     cfg.Public.GenesisBlock.CoinVolume,
			BodyHash: bodyHash,
			PrevHash: cipher.SHA256{},
			BkSeq:    0,
			Version:  0,
			Fee:      0,
			UxHash:   cipher.SHA256{},
		},

		Body: coin.BlockBody{
			Transactions: coin.Transactions{tx},
		},
	}

	sb := coin.SignedBlock{
		Block: b,
		Sig:   sig,
	}

	return sb, nil
}

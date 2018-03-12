package main

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// Config records the node's configuration
type NodeConfig struct {
	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Disables wallet API
	DisableWalletApi bool
	// Disable CSRF check in the wallet api
	DisableCSRF bool

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max outgoing connections to maintain
	MaxOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
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

	GenesisSignature  cipher.Sig
	GenesisAddress    cipher.Address
	GenesisTimestamp  uint64
	GenesisCoinVolume uint64

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

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
	RPCThreadNum uint // rpc number
	Logtofile    bool
	Logtogui     bool
	LogBuffSize  int
}

func makeDefaultNodeConfig() NodeConfig {
	cfg := NodeConfig{
		// Disable peer exchange
		DisablePEX: false,
		// Don't make any outgoing connections
		DisableOutgoingConnections: false,
		// Don't allowing incoming connections
		DisableIncomingConnections: false,
		// Disables networking altogether
		DisableNetworking: false,
		// Disable wallet API
		DisableWalletApi: false,
		// Disable CSRF check in the wallet api
		DisableCSRF: true,
		// Only run on localhost and only connect to others on localhost
		LocalhostOnly: false,
		// Which address to serve on. Leave blank to automatically assign to a
		// public interface
		Address: "",
		//gnet uses this for TCP incoming and outgoing
		Port: 16000,
		// MaxOutgoingConnections is the maximum outgoing connections allowed.
		MaxOutgoingConnections: 16,
		DownloadPeerList:       false,
		PeerListURL:            "https://downloads.skycoin.net/blockchain/peers.txt",
		// How often to make outgoing connections, in seconds
		OutgoingConnectionsRate: time.Second * 5,
		PeerlistSize:            65535,
		// Wallet Address Version
		//AddressVersion: "test",
		// Remote web interface
		WebInterface:             true,
		WebInterfacePort:         6420,
		WebInterfaceAddr:         "127.0.0.1",
		WebInterfaceCert:         "",
		WebInterfaceKey:          "",
		WebInterfaceHTTPS:        false,
		PrintWebInterfaceAddress: false,

		RPCInterface:     true,
		RPCInterfacePort: 6430,
		RPCInterfaceAddr: "127.0.0.1",
		RPCThreadNum:     5,

		LaunchBrowser: true,
		// Data directory holds app data -- defaults to ~/.skycoin
		DataDirectory: ".skycoin",
		// Web GUI static resources
		GUIDirectory: "./src/gui/static/",
		// Logging
		ColorLog: true,
		LogLevel: "DEBUG",

		// Wallets
		WalletDirectory: "",

		/* Developer options */

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "skycoin.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf: false,
		// Will force it to connect to this ip:port, instead of waiting for it
		// to show up as a peer
		ConnectTo:   "",
		LogBuffSize: 8388608, //1024*1024*8
	}

	return cfg
}

func makeNodeConfig(toolCfg common.Config) (NodeConfig, error) {
	var (
		cfg = makeDefaultNodeConfig()
		err error
	)

	// Master's key par
	if cfg.BlockchainSeckey, err = cipher.SecKeyFromHex(toolCfg.Secret.MasterSecKey); err != nil {
		return cfg, errors.New("invalid master node secret key")
	}
	if cfg.BlockchainPubkey, err = cipher.PubKeyFromHex(toolCfg.Public.MasterPubKey); err != nil {
		return cfg, errors.New("invalid master node public key")
	}

	// Genesis block
	gbCfg := toolCfg.Public.GenesisBlock

	if cfg.GenesisSignature, err = cipher.SigFromHex(toolCfg.Secret.GenesisSignature); err != nil {
		return cfg, errors.New("invalid genesis signature")
	}
	if cfg.GenesisAddress, err = cipher.DecodeBase58Address(gbCfg.Address); err != nil {
		return cfg, errors.New("invalid genesis address")
	}
	cfg.GenesisCoinVolume = gbCfg.CoinVolume * 1e6
	cfg.GenesisTimestamp = gbCfg.Timestamp

	// Network
	cfg.Port = toolCfg.Public.Port
	cfg.RPCInterfacePort = toolCfg.Public.RPCPort
	cfg.WebInterfacePort = toolCfg.Public.GUIPort

	// Data directory
	if _, err = file.InitDataDir(cfg.DataDirectory); err != nil {
		return cfg, err
	}
	cfg.WebInterfaceCert = filepath.Join(cfg.DataDirectory, "cert.pem")
	cfg.WebInterfaceKey = filepath.Join(cfg.DataDirectory, "key.pem")
	cfg.WalletDirectory = filepath.Join(cfg.DataDirectory, "wallets")
	cfg.DBPath = filepath.Join(cfg.DataDirectory, "data.db")

	cfg.GUIDirectory = file.ResolveResourceDirectory(cfg.GUIDirectory)

	return cfg, nil
}

func makeDaemonConfg(nc NodeConfig) daemon.Config {
	dc := daemon.NewConfig()

	// PEX
	dc.Pex.DataDirectory = nc.DataDirectory
	dc.Pex.Disabled = nc.DisablePEX
	dc.Pex.Max = nc.PeerlistSize
	dc.Pex.DownloadPeerList = nc.DownloadPeerList
	dc.Pex.PeerListURL = nc.PeerListURL
	dc.Pex.AllowLocalhost = true

	// Networking
	dc.Daemon.DisableOutgoingConnections = nc.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = nc.DisableIncomingConnections
	dc.Daemon.DisableNetworking = nc.DisableNetworking
	dc.Daemon.Port = nc.Port
	dc.Daemon.Address = nc.Address
	dc.Daemon.LocalhostOnly = nc.LocalhostOnly
	dc.Daemon.OutgoingMax = nc.MaxOutgoingConnections
	dc.Daemon.DataDirectory = nc.DataDirectory
	dc.Daemon.LogPings = !nc.DisablePingPong
	dc.Daemon.OutgoingRate = nc.OutgoingConnectionsRate

	// Centralized network configuration
	dc.Visor.Config.IsMaster = nc.RunMaster
	dc.Visor.Config.BlockchainPubkey = nc.BlockchainPubkey
	dc.Visor.Config.BlockchainSeckey = nc.BlockchainSeckey
	dc.Visor.Config.GenesisSignature = nc.GenesisSignature
	dc.Visor.Config.GenesisAddress = nc.GenesisAddress
	dc.Visor.Config.GenesisCoinVolume = nc.GenesisCoinVolume
	dc.Visor.Config.GenesisTimestamp = nc.GenesisTimestamp

	// Wallet
	dc.Visor.Config.WalletDirectory = nc.WalletDirectory
	dc.Gateway.DisableWalletAPI = nc.DisableWalletApi

	dc.Visor.Config.DBPath = nc.DBPath

	dc.Visor.Config.Arbitrating = nc.Arbitrating

	dc.Visor.Config.BuildInfo = visor.BuildInfo{
		Version: Version,
		Commit:  Commit,
	}

	return dc
}

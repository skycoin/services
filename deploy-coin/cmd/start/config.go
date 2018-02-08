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

	GenesisSignature  cipher.Sig
	GenesisTimestamp  uint64
	GenesisCoinVolume uint64
	GenesisAddress    cipher.Address

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

func makeNodeConfig(toolCfg common.Config, runMaster bool) (NodeConfig, error) {
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

	// Master node is the only trusted peer of new network
	cfg.RunMaster = runMaster
	cfg.LocalhostOnly = true
	/*
		if !cfg.RunMaster {
			cfg.DefaultConnections = []string{
				fmt.Sprintf("127.0.0.1:%d", cfg.Port),
			}
		}
	*/

	// Network
	cfg.Port = toolCfg.Public.Port
	cfg.WebInterfacePort = toolCfg.Public.WebInterfacePort
	cfg.RPCInterfacePort = toolCfg.Public.RPCInterfacePort

	// Data directory
	cfg.DataDirectory = toolCfg.Public.DataDirectory
	if _, err = file.InitDataDir(cfg.DataDirectory); err != nil {
		return cfg, err
	}

	cfg.WebInterfaceCert = filepath.Join(cfg.DataDirectory, "cert.pem")
	cfg.WebInterfaceKey = filepath.Join(cfg.DataDirectory, "key.pem")
	cfg.WalletDirectory = filepath.Join(cfg.DataDirectory, "wallets")
	cfg.DBPath = filepath.Join(cfg.DataDirectory, "data.db")

	cfg.LogFmt = toolCfg.Public.LogFmt

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
	cfg.GenesisCoinVolume = gbCfg.CoinVolume
	cfg.GenesisTimestamp = gbCfg.Timestamp

	return cfg, nil
}

func makeDaemonConfg(nc NodeConfig) daemon.Config {
	dc := daemon.NewConfig()

	dc.Peers.DataDirectory = nc.DataDirectory
	dc.Peers.Disabled = nc.DisablePEX
	dc.Daemon.DisableOutgoingConnections = nc.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = nc.DisableIncomingConnections
	dc.Daemon.DisableNetworking = nc.DisableNetworking
	dc.Daemon.Port = nc.Port
	dc.Daemon.Address = nc.Address
	dc.Daemon.LocalhostOnly = nc.LocalhostOnly
	dc.Daemon.OutgoingMax = nc.MaxConnections
	dc.Daemon.DataDirectory = nc.DataDirectory
	dc.Daemon.LogPings = !nc.DisablePingPong

	daemon.DefaultConnections = nc.DefaultConnections

	dc.Daemon.OutgoingRate = nc.OutgoingConnectionsRate

	dc.Visor.Config.IsMaster = nc.RunMaster

	dc.Visor.Config.BlockchainPubkey = nc.BlockchainPubkey
	dc.Visor.Config.BlockchainSeckey = nc.BlockchainSeckey

	dc.Visor.Config.GenesisSignature = nc.GenesisSignature
	dc.Visor.Config.GenesisAddress = nc.GenesisAddress
	dc.Visor.Config.GenesisCoinVolume = nc.GenesisCoinVolume
	dc.Visor.Config.GenesisTimestamp = nc.GenesisTimestamp

	dc.Visor.Config.DBPath = nc.DBPath
	dc.Visor.Config.Arbitrating = nc.Arbitrating
	dc.Visor.Config.WalletDirectory = nc.WalletDirectory
	dc.Visor.Config.BuildInfo = visor.BuildInfo{
		Version: Version,
		Commit:  Commit,
	}

	return dc
}

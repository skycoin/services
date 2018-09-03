package common

type Config struct {
	Secret SecretConfig `json:"secret"`
	Public PublicConfig `json:"public"`
}

type SecretConfig struct {
	MasterSecKey     string `json:"masterPrivateKey"`
	GenesisSignature string `json:"genesisSignature"`
}

type PublicConfig struct {
	MasterPubKey string `json:"masterPublicKey"`

	GenesisBlock  GenesisBlockConfig  `json:"genesisBlock"`
	GenesisWallet GenesisWalletConfig `json:"genesisWallet"`

	CoinCode string `json:"coinCode"`

	Port    int `json:"port"`
	RPCPort int `json:"rpcPort"`
	GUIPort int `json:"guiPort"`

	TrustedPeers []string `json:"trustedPeers"`
}

type GenesisBlockConfig struct {
	Address    string `json:"address"`
	CoinVolume uint64 `json:"coins"`
	Timestamp  uint64 `json:"timestamp"`
	BodyHash   string `json:"bodyHash"`
	HeaderHash string `json:"headerHash"`
}

type GenesisWalletConfig struct {
	Seed            string `json:"seed"`
	Addresses       uint64 `json:"addresses"`
	CoinsPerAddress uint64 `json:"coinsPerAddress"`
}

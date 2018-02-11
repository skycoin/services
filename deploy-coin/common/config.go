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

	GenesisBlock GenesisBlockConfig `json:"genesisBlock"`

	CoinCode string `json:"coinCode"`

	Distribution DistributionConfig `json:"distribuion"`

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

type DistributionConfig struct {
	CoinsPerAddress uint64   `json:"coinsPerAddress"`
	AddressSeed     string   `json:"addressSeed"`
	Addresses       []string `json:"addresses"`
}

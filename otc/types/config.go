package types

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Skycoin struct {
		Node string
		Seed string
		Name string
	}
	CoinApi struct {
		URL string
	}
	Dropper struct {
		BTC struct {
			Node    string
			User    string
			Pass    string
			Account string
			Price   uint64
		}
	}
	Api struct {
		Listen string
	}
	Admin struct {
		Listen string
	}
	Model struct {
		Tick   int
		Path   string
		Paused bool
	}
	Scanner struct {
		Tick       int
		Expiration int
	}
	Sender struct {
		Tick int
	}
	Monitor struct {
		Tick int
	}
}

func NewConfig(path string) (*Config, error) {
	c := &Config{}
	_, err := toml.DecodeFile(path, &c)
	return c, err
}

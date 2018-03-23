package otc

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	SKY struct {
		Node string
		Seed string
		Name string
	}
	BTC struct {
		Node    string
		User    string
		Pass    string
		Account string
		Testnet bool
	}
	API struct {
		Public struct {
			Listen string
		}
		Admin struct {
			Listen string
		}
	}
	Watcher struct {
		Node string
	}
}

func NewConfig(path string) (*Config, error) {
	c := &Config{}
	_, err := toml.DecodeFile(path, &c)
	return c, err
}

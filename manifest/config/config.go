package config

import (
	"encoding/json"
	"fmt"
	"os"
)

//Configuration
type Configuration struct {
	Folders     []string `json:"folders"`
}

//Config is a global variable
var Config Configuration

//LoadConfiguration setup config
func LoadConfiguration(file string) Configuration {

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Config)
	return Config
}

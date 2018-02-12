package main

import (
	"github.com/skycoin/services/manifest/config"
	"github.com/skycoin/services/manifest/manifest"
	"fmt"
)



func main() {


	config.LoadConfiguration("./config.json")
	fmt.Println(config.Config)
	manifest.ReadFiles(config.Config.Folders)


}




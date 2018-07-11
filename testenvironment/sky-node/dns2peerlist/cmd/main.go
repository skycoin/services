package main

import (
	"flag"
	"net"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/daemon/pex"
	"time"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"bufio"
)

type Config struct {
	port string
	outputFile string
	service string
	format string
}

func newConfig() *Config {
	c := &Config{}
	flag.StringVar(&c.port,"port","6000","port to set to every peer on the peerlist")
	flag.StringVar(&c.outputFile,"output","./peers.txt","path to output file")
	flag.StringVar(&c.service,"service","skycoin-node","name of the service to ")
	flag.StringVar(&c.format,"format","json","which format to write the output in: json or text")

	flag.Parse()

	return c
	}

func main() {
	c := newConfig()

	ips, err := net.LookupIP(c.service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not resolve %s\n", c.service)
		os.Exit(1)
	}

	switch c.format {
	case "json":
		jsonWrite(c,ips)
	case "text":
		textWrite(c,ips)
	default:
		fmt.Println("format not allowed, using default value")
		jsonWrite(c,ips)
	}


}


func jsonWrite(c *Config, ips []net.IP) {
	peerList := make(map[string]pex.PeerJSON)

	for _, ip := range ips {
		addrs := ip.String()+":"+c.port
		peerList[addrs]=peer(addrs)
	}

	encodedPeerlist, err := json.Marshal(peerList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse to json: %+v\n", peerList)
		os.Exit(1)
	}

	prettyPeerList, _ := prettyprint(encodedPeerlist)
	err = ioutil.WriteFile(c.outputFile,prettyPeerList, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write file at %s\n", c.outputFile)
		os.Exit(1)
	}
}

func textWrite(c *Config, ips []net.IP) {
	peerList := []string{}
	for _, ip := range ips {
		addrs := ip.String()+":"+c.port
		peerList = append(peerList, addrs)
	}

	err := writeLines(peerList, c.outputFile)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Could not write file at %s\n", c.outputFile)
		os.Exit(1)
	}
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func peer(addrs string) pex.PeerJSON {
	incomingPort := true

	return pex.PeerJSON{
		Addr: addrs,
		Private: true,
		Trusted: true,
		LastSeen: time.Now().Unix(),
		HasIncomingPort: &incomingPort,
	}
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
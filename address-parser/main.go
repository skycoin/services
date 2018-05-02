package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	expand "github.com/openvenues/gopostal/expand"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

type Item struct {
	Original string   `json:"original"`
	Parsed   []string `json:"parsed"`
}

func getCommands() map[string]*cobra.Command {
	return map[string]*cobra.Command{
		"file": {
			Use:     "file",
			Short:   "",
			Example: "address-parser file <path-to-file>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				processFromFile(args[0])
			},
		},
		"address": {
			Use:     "address",
			Short:   "",
			Example: "address-parser address <address-string>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				processFromAddress(args[0])
			},
		},
	}
}

func init() {
	rootCmd = &cobra.Command{Use: "address-parser"}
	for _, command := range getCommands() {
		rootCmd.AddCommand(command)
	}
}

func processFromAddress(addr string) {
	var items = make([]*Item, 0)
	expansions := expand.ExpandAddress(addr)
	item := &Item{
		Original: addr,
		Parsed:   expansions,
	}
	items = append(items, item)
	result, err := json.MarshalIndent(&items, "", "    ")
	if err != nil {
		fmt.Printf("failed get address. err: %v\n", err.Error())
	}
	fmt.Println(string(result))
}

func processFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("failed open file: %v. err: %v\n", path, err.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var items = make([]*Item, 0)

	for scanner.Scan() {
		addr := scanner.Text()
		if addr != "" {
			expansions := expand.ExpandAddress(addr)
			item := &Item{
				Original: addr,
				Parsed:   expansions,
			}
			items = append(items, item)
		}
	}
	result, err := json.MarshalIndent(&items, "", "    ")
	if err != nil {
		fmt.Printf("failed get address. err: %v\n", err.Error())
		return
	}
	fmt.Println(string(result))
}

func main() {
	rootCmd.Execute()
}

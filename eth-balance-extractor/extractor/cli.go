package extractor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	gcli "github.com/urfave/cli"
)

// App represents CLI app
type App struct {
	gcli.App
}

// NewApp creates a new instance of the App
func NewApp() *App {
	gcliApp := gcli.NewApp()

	commands := []gcli.Command{
		extractWallets(),
		continueExtraction(),
	}

	gcliApp.Commands = commands

	return &App{
		App: *gcliApp,
	}
}

// Run starts the app
func (app *App) Run(args []string) error {
	return app.App.Run(args)
}

func onCommandUsageError(command string) gcli.OnUsageErrorFunc {
	return func(c *gcli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "Error: %v\n\n", err)
		return gcli.ShowCommandHelp(c, command)
	}
}

func extractWallets() gcli.Command {
	name := "extractWallets"
	return gcli.Command{
		Name:         name,
		Usage:        "Starts extraction process",
		ArgsUsage:    "[node_api_url] [smart_contract_address] [smart_contract_transfer_method_hash] [smart_contract_transfer_from_method_hash] [dest_dir] [start_block] [threads_count]",
		Description:  fmt.Sprintf(`Starts extraction process`),
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			nodeAPIUrl := c.Args().Get(0)
			smartContractAddress := c.Args().Get(1)
			transferHash := c.Args().Get(2)
			transferFromHash := c.Args().Get(3)
			destDir := c.Args().Get(4)
			startBlock, err := strconv.Atoi(c.Args().Get(5))
			if err != nil {
				fmt.Println("cli > ", err)
				return gcli.ShowSubcommandHelp(c)
			}
			threadsCount, err := strconv.Atoi(c.Args().Get(6))
			if err != nil {
				fmt.Println("cli > ", err)
				return gcli.ShowSubcommandHelp(c)
			}

			o := NewOrchestrator(nodeAPIUrl, smartContractAddress, transferHash, transferFromHash, destDir, startBlock, threadsCount)

			go o.StartScanning()

			reader := bufio.NewReader(os.Stdin)
			for {
				reader.ReadString('\n')
			}

			return nil
		},
	}
}

func continueExtraction() gcli.Command {
	name := "continueExtraction"
	return gcli.Command{
		Name:         name,
		Usage:        "Continue extraction process",
		ArgsUsage:    "[wallets_file] [node_api_url] [smart_contract_address] [smart_contract_transfer_method_hash] [smart_contract_transfer_from_method_hash] [dest_dir] [start_block] [threads_count]",
		Description:  fmt.Sprintf(`Starts extraction process`),
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			walletsFile := c.Args().Get(0)
			nodeAPIUrl := c.Args().Get(1)
			smartContractAddress := c.Args().Get(2)
			transferHash := c.Args().Get(3)
			transferFromHash := c.Args().Get(4)
			destDir := c.Args().Get(5)
			startBlock, err := strconv.Atoi(c.Args().Get(6))
			if err != nil {
				fmt.Println("cli > ", err)
				return gcli.ShowSubcommandHelp(c)
			}
			threadsCount, err := strconv.Atoi(c.Args().Get(7))
			if err != nil {
				fmt.Println("cli > ", err)
				return gcli.ShowSubcommandHelp(c)
			}

			storage := NewStorage(destDir)
			wallets := storage.LoadSnapshot(walletsFile)
			o := NewOrchestrator(nodeAPIUrl, smartContractAddress, transferHash, transferFromHash, destDir, startBlock, threadsCount)
			o.scanner.Wallets = wallets

			go o.StartScanning()

			reader := bufio.NewReader(os.Stdin)
			for {
				reader.ReadString('\n')
			}

			return nil
		},
	}
}

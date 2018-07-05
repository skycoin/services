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
		ArgsUsage:    "[node_api_url] [smart_contract_address] [dest_dir] [start_block] [threads_count]",
		Description:  fmt.Sprintf(`Starts extraction process`),
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			nodeAPIUrl := c.Args().Get(0)
			smartContractAddress := c.Args().Get(1)
			destDir := c.Args().Get(2)
			startBlock, err := strconv.Atoi(c.Args().Get(3))
			if err != nil {
				fmt.Println(err)
				return gcli.ShowSubcommandHelp(c)
			}
			threadsCount, err := strconv.Atoi(c.Args().Get(4))
			if err != nil {
				fmt.Println(err)
				return gcli.ShowSubcommandHelp(c)
			}

			o := NewOrchestrator(nodeAPIUrl, smartContractAddress, destDir, startBlock, threadsCount)

			go o.StartScanning()

			reader := bufio.NewReader(os.Stdin)
			for {
				reader.ReadString('\n')
			}

			return nil
		},
	}
}

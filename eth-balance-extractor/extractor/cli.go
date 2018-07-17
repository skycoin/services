package extractor

import (
	"fmt"

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
		extractWalletsPublicKeys(),
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

func extractWalletsPublicKeys() gcli.Command {
	name := "extractWalletsKeys"
	return gcli.Command{
		Name:         name,
		Usage:        "Extracts wallets public keys by provided tx hashes",
		ArgsUsage:    "[node_api_url] [wallets_file] [dest_dir]",
		Description:  fmt.Sprintf(`Extracts wallets public keys by provided tx hashes`),
		OnUsageError: onCommandUsageError(name),

		Action: func(c *gcli.Context) error {
			nodeAPIUrl := c.Args().Get(0)
			walletsFile := c.Args().Get(1)
			destDir := c.Args().Get(2)

			s := NewStorage(destDir)
			wallets := s.LoadTransactionWallets(walletsFile)

			scanner := NewWalletScanner(nodeAPIUrl, wallets)
			scanner.RestoreKeys()

			s.StoreSnapshot("wallets", scanner.Wallets)

			return nil
		},
	}
}

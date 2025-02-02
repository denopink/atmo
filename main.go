package main

import (
	"github.com/spf13/cobra"

	"github.com/suborbital/e2core/e2core/command"
	"github.com/suborbital/e2core/e2core/release"
)

func main() {
	root := rootCommand()
	root.AddCommand(command.Start())

	mod := modCommand()
	mod.AddCommand(command.ModStart())
	root.AddCommand(mod)

	root.Execute()
}

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "e2core",
		Version: release.Version,
		Long: `
	E2Core is a secure development kit and server for writing and running untrusted third-party plugins.

	The E2Core server is responsible for managing and running plugins using simple HTTP, RPC, or streaming interfaces.`,
	}

	return cmd
}

func modCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mod",
		Version: release.Version,
		Short:   "commands for working with modules",
		Hidden:  true,
	}

	return cmd
}

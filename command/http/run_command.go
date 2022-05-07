package http

import "github.com/spf13/cobra"

func init() {
	command.AddCommand(runCommand)
}

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "",
	Args:  cobra.NoArgs,
	Run:   run,
}

func run(*cobra.Command, []string) {}

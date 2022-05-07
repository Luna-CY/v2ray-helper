package http

import "github.com/spf13/cobra"

func init() {
	command.Flags().StringVarP(&home, "home", "d", "", "工作目录")
	command.Flags().StringVarP(&env, "env", "e", "prod", "运行环境")
}

var home string
var env string

var command = &cobra.Command{
	Use:   "http-service",
	Short: "",
	Args:  cobra.NoArgs,
}

func Execute() error {
	return command.Execute()
}

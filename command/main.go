package command

import (
	"github.com/spf13/cobra"
	"os"
)

var command = cobra.Command{
	Use:   "v2ray-helper",
	Short: "V2ray配置服务，提供对V2ray的可视化配置操作",
}

var (
	home string
	host string
)

func Exec() {
	if err := command.Execute(); nil != err {
		os.Exit(1)
	}
}

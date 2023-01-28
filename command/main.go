package command

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	home    string
	host    string
	command cobra.Command
)

func init() {
	command = cobra.Command{
		Use:              "v2ray-helper",
		Short:            "V2ray配置服务，提供对V2ray的可视化配置操作",
		PersistentPreRun: initialize,
	}

	command.PersistentFlags().StringVar(&home, "home", "", "运行主目录，默认为服务命令所在目录")

	command.AddCommand(runCommand)
	command.AddCommand(certCommand)
	command.AddCommand(testCommand)
	command.AddCommand(installCommand)
}

func Exec() {
	if err := command.Execute(); nil != err {
		os.Exit(1)
	}
}

// initialize 初始化系统的公共方法
// 使用此方法的命令必须添加 home 变量的参数解析
// 此方法提供对配置、日志、证书管理器的初始化
func initialize(cmd *cobra.Command, _ []string) {
	time.Local = time.FixedZone("CST", 8*3600)
	home = filepath.Clean(strings.TrimSpace(home))
	rootAbsPath := runtime.AbsRootPath(home)

	if err := configurator.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	// 设置GIN模式需要在其他组件初始化之前
	if configurator.Configure.Server.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	// logger组件需要在其他组件之前初始化
	if err := logger.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化日志失败: %v", err))
	}

	if err := certificate.Init(cmd.Context()); nil != err {
		log.Fatalln(fmt.Sprintf("初始化证书管理器失败: %v", err))
	}
}

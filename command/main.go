package command

import (
	"context"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		PersistentPreRun: initSystem,
	}

	command.PersistentFlags().StringVar(&home, "home", "", "运行主目录，默认为服务命令所在目录")
}

func Exec() {
	if err := command.Execute(); nil != err {
		os.Exit(1)
	}
}

// initSystem 初始化系统的公共方法
// 使用此方法的命令必须添加 home 变量的参数解析
// 此方法提供对配置、日志、证书管理器的初始化
func initSystem(*cobra.Command, []string) {
	home = filepath.Clean(strings.TrimSpace(home))
	rootAbsPath := runtime.AbsRootPath(home)

	if err := configurator.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	if err := runtime.InitRuntime(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化运行环境失败: %v", err))
	}

	// 设置GIN模式需要在其他组件初始化之前
	if viper.GetBool(configurator.KeyServerRelease) {
		gin.SetMode(gin.ReleaseMode)
	}

	// logger组件需要在其他组件之前初始化
	if err := logger.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化日志失败: %v", err))
	}

	if err := certificate.Init(context.Background()); nil != err {
		log.Fatalln(fmt.Sprintf("初始化证书管理器失败: %v", err))
	}

	if err := database.Init(filepath.Join(rootAbsPath, "main.db"), 10); nil != err {
		log.Fatalln(fmt.Sprintf("初始化数据库失败: %v", err))
	}
}

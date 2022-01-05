package configurator

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	DefaultAccessKey     = "-"
	DefaultManagementKey = "-"
)

const (
	KeyRootPath          = "root-path"
	KeyServerAddress     = "server.address"
	KeyServerPort        = "server.port"
	KeyServerHttpsEnable = "server.https.enable"
	KeyServerHttpsHost   = "server.https.host"
	KeyServerRelease     = "server.release"
	KeyServerAllowDeploy = "server.allow-deploy"
	KeyHttpsIssueEmail   = "https.issue-email"
	KeyAuthAccessKey     = "auth.access-key"
	KeyAuthManagementKey = "auth.management-key"
	KeyLogLevel          = "log.level"
)

func Init(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); nil != err {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("无法读取配置文件: %v\n", err)

			os.Exit(1)
		}

		log.Println(err)
		log.Println("使用默认配置启动")
	}

	viper.SetDefault(KeyRootPath, path)
	viper.SetDefault(KeyServerAddress, "0.0.0.0")
	viper.SetDefault(KeyServerPort, 8888)
	viper.SetDefault(KeyServerHttpsEnable, false)
	viper.SetDefault(KeyServerHttpsHost, "")
	viper.SetDefault(KeyServerRelease, true)
	viper.SetDefault(KeyServerAllowDeploy, true)
	viper.SetDefault(KeyAuthAccessKey, DefaultAccessKey)
	viper.SetDefault(KeyAuthManagementKey, DefaultManagementKey)
	viper.SetDefault(KeyHttpsIssueEmail, fmt.Sprintf("%v@v2ray-helper.net", util.GenerateRandomString(16)))
	viper.SetDefault(KeyLogLevel, "error")

	return nil
}

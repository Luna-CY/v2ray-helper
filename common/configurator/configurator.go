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
	KeyRootPath              = "root-path"
	KeyServerAddress         = "server.address"
	KeyServerPort            = "server.port"
	KeyServerRelease         = "server.release"
	KeyAcmeEmail             = "acme.email"
	KeyAuthAccessKey         = "auth.access"
	KeyAuthManagementKey     = "auth.manager"
	KeyMailSMTPServerAddress = "mail.smtp.address"
	KeyMailSMTPServerPort    = "mail.smtp.port"
	KeyMailSMTPUser          = "mail.smtp.user"
	KeyMailSMTPPassword      = "mail.smtp.password"
	KeyMailSMTPSecret        = "mail.smtp.secret"
	KeyMailNoticeEnable      = "mail.notice.enable"
	KeyMailNoticeTo          = "mail.notice.to"
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
	viper.SetDefault(KeyServerAddress, "127.0.0.1")
	viper.SetDefault(KeyServerPort, 8888)
	viper.SetDefault(KeyServerRelease, true)
	viper.SetDefault(KeyAuthAccessKey, DefaultAccessKey)
	viper.SetDefault(KeyAuthManagementKey, DefaultManagementKey)
	viper.SetDefault(KeyAcmeEmail, fmt.Sprintf("%v@v2ray-helper.net", util.GenerateRandomString(16)))
	viper.SetDefault(KeyMailSMTPServerAddress, "")
	viper.SetDefault(KeyMailSMTPServerPort, 587)
	viper.SetDefault(KeyMailSMTPUser, "")
	viper.SetDefault(KeyMailSMTPPassword, "")
	viper.SetDefault(KeyMailSMTPSecret, "")
	viper.SetDefault(KeyMailNoticeEnable, false)
	viper.SetDefault(KeyMailNoticeTo, "")

	Configure.Home = path
	if err := viper.Unmarshal(&Configure); nil != err {
		return err
	}

	return nil
}

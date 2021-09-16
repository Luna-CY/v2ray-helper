package configurator

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

const DefaultMainConfigContent = `address: 0.0.0.0
service-listen: 8888
https-listen: 8888
gin-release-mode: true
key: '-'
remove-key: '-'
email: myself@v2ray-helper.net
allow-v2ray-deploy: true
log-level: error`

const DefaultMainConfigWithHttpsContent = `address: 127.0.0.1
service-listen: 9999
https-listen: 8888
gin-release-mode: true
key: '-'
remove-key: '-'
email: myself@v2ray-helper.net
allow-v2ray-deploy: true
log-level: error`

type mainConfig struct {
	Address          string `yaml:"address"`
	ServiceListen    int    `yaml:"service-listen"`
	HttpsListen      int    `yaml:"https-listen"`
	GinReleaseMode   bool   `yaml:"gin-release-mode"`
	AllowV2rayDeploy bool   `yaml:"allow-v2ray-deploy"`
	Email            string `yaml:"email"`
	Key              string `yaml:"key"`
	RemoveKey        string `yaml:"remove-key"`
	LogLevel         string `yaml:"log-level"`
}

func (m *mainConfig) GetListenAddress() string {
	return fmt.Sprintf("%v:%v", m.Address, m.ServiceListen)
}

func (m *mainConfig) GetLogLevel() logrus.Level {
	maps := map[string]logrus.Level{"debug": logrus.DebugLevel, "info": logrus.InfoLevel, "warn": logrus.WarnLevel, "error": logrus.ErrorLevel}
	if level, ok := maps[m.LogLevel]; ok {
		return level
	}

	return logrus.ErrorLevel
}

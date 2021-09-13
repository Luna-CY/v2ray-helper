package configurator

import (
	"github.com/sirupsen/logrus"
	"strings"
)

const DefaultMainConfigContent = `listen: 0.0.0.0:8888
key: '-'
remove-key: '-'
email: myself@v2ray-helper.net
allow-v2ray-deploy: true
log-level: error`

type mainConfig struct {
	Listen           string `yaml:"listen"`
	AllowV2rayDeploy bool   `yaml:"allow-v2ray-deploy"`
	Email            string `yaml:"email"`
	Key              string `yaml:"key"`
	RemoveKey        string `yaml:"remove-key"`
	LogLevel         string `yaml:"log-level"`
}

func (m *mainConfig) GetListen() string {
	if !strings.Contains(m.Listen, ":") {
		return "0.0.0.0:8888"
	}

	return strings.TrimSpace(m.Listen)
}

func (m *mainConfig) GetLogLevel() logrus.Level {
	maps := map[string]logrus.Level{"debug": logrus.DebugLevel, "info": logrus.InfoLevel, "warn": logrus.WarnLevel, "error": logrus.ErrorLevel}
	if level, ok := maps[m.LogLevel]; ok {
		return level
	}

	return logrus.ErrorLevel
}

package configurator

import (
	"github.com/sirupsen/logrus"
	"strings"
)

const defaultLogPath = "/var/log/v2ray-subscription.log"

type mainConfig struct {
	Listen    int    `yaml:"listen"`
	Key       string `yaml:"key"`
	RemoveKey string `yaml:"remove-key"`
	LogLevel  string `yaml:"log-level"`
	LogPath   string `yaml:"log-path"`
}

func (m *mainConfig) GetLogPath() string {
	m.LogPath = strings.TrimSpace(m.LogPath)

	if "" == m.LogPath {
		return defaultLogPath
	}

	return m.LogPath
}

func (m *mainConfig) GetLogLevel() logrus.Level {
	maps := map[string]logrus.Level{"debug": logrus.DebugLevel, "info": logrus.InfoLevel, "warn": logrus.WarnLevel, "error": logrus.ErrorLevel}
	if level, ok := maps[m.LogLevel]; ok {
		return level
	}

	return logrus.ErrorLevel
}

package configurator

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type mainConfig struct {
	Address          string `yaml:"address"`
	Listen           int    `yaml:"listen"`
	EnableHttps      bool   `yaml:"enable-https"`
	HttpsHost        string `yaml:"https_host"`
	GinReleaseMode   bool   `yaml:"gin-release-mode"`
	AllowV2rayDeploy bool   `yaml:"allow-v2ray-deploy"`
	Email            string `yaml:"email"`
	AccessKey        string `yaml:"access-key"`
	ManagementKey    string `yaml:"management-key"`
	LogLevel         string `yaml:"log-level"`
}

func (m *mainConfig) GetListenAddress() string {
	return fmt.Sprintf("%v:%v", m.Address, m.Listen)
}

func (m *mainConfig) GetLogLevel() logrus.Level {
	maps := map[string]logrus.Level{"debug": logrus.DebugLevel, "info": logrus.InfoLevel, "warn": logrus.WarnLevel, "error": logrus.ErrorLevel}
	if level, ok := maps[m.LogLevel]; ok {
		return level
	}

	return logrus.ErrorLevel
}

func (m *mainConfig) GetFileName() string {
	return "main.config.yaml"
}

// Load 加载配置
func (m *mainConfig) Load(configPath string) error {
	configFile, err := os.Open(configPath)
	if nil != err {
		return errors.New(fmt.Sprintf("找到不配置文件: %v %v", configPath, err))
	}
	defer configFile.Close()

	configContent, err := ioutil.ReadAll(configFile)
	if nil != err {
		return errors.New(fmt.Sprintf("无法读取配置文件: %v %v", configPath, err))
	}

	if err := yaml.Unmarshal(configContent, m); nil != err {
		return errors.New(fmt.Sprintf("无法解析配置文件: %v %v", configPath, err))
	}
	return nil
}

// Save 保存配置到文件
func (m *mainConfig) Save(configPath string) error {
	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开配置文件: %v %v", configPath, err))
	}
	defer configFile.Close()

	content, err := yaml.Marshal(m)
	if nil != err {
		return errors.New(fmt.Sprintf("无法序列化配置参数: %v", err))
	}

	if _, err := configFile.Write(content); nil != err {
		return errors.New(fmt.Sprintf("无法写入配置文件: %v", err))
	}

	return nil
}

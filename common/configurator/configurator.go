package configurator

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	mc *mainConfig
)

const (
	DefaultKey       = "-"
	DefaultRemoveKey = "-"
)

func Init(rootPath string) error {

	mc = &mainConfig{
		Listen:           "0.0.0.0:8888",
		GinReleaseMode:   true,
		Email:            "myself@v2ray-helper.net",
		AllowV2rayDeploy: true,
		Key:              DefaultKey,
		RemoveKey:        DefaultRemoveKey,
		LogLevel:         "error",
	}

	if err := loadConfig(filepath.Join(rootPath, "config", "main.prod.config.yaml"), filepath.Join(rootPath, "config", "main.local.config.yaml"), &mc); nil != err {
		return err
	}

	return nil
}

func GetMainConfig() *mainConfig {
	if nil == mc {
		panic("获取了未初始化的配置")
	}

	return mc
}

func loadConfig(configPath, localPath string, dest interface{}) error {
	configFile, err := os.Open(configPath)
	if nil != err {
		return errors.New(fmt.Sprintf("找到不配置文件: %v %v", configPath, err))
	}
	defer configFile.Close()

	configContent, err := ioutil.ReadAll(configFile)
	if nil != err {
		return errors.New(fmt.Sprintf("无法读取配置文件: %v %v", configPath, err))
	}

	if err := yaml.Unmarshal(configContent, dest); nil != err {
		return errors.New(fmt.Sprintf("无法解析配置文件: %v %v", configPath, err))
	}

	if _, err = os.Stat(localPath); nil == err {
		localFile, err := os.Open(localPath)
		if nil != err {
			return errors.New(fmt.Sprintf("找到不配置文件: %v %v", localPath, err))
		}
		defer localFile.Close()

		localConfigContent, err := ioutil.ReadAll(localFile)
		if nil != err {
			return errors.New(fmt.Sprintf("无法读取配置文件: %v %v", localPath, err))
		}

		if err := yaml.Unmarshal(localConfigContent, dest); nil != err {
			return errors.New(fmt.Sprintf("无法解析配置文件: %v %v", localPath, err))
		}
	}

	return nil
}

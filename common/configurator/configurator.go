package configurator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	mc *mainConfig
	dc *dbConfig
)

func Init() error {
	var rootPath string

	if gin.ReleaseMode == gin.Mode() {
		rootPath = getRootPath()
	} else {
		pwd, err := os.Getwd()
		if nil != err {
			return errors.New(fmt.Sprintf("无法获取当前目录: %v", err))
		}

		rootPath = path.Dir(pwd)
	}

	mc = &mainConfig{
		Listen:             "0.0.0.0:8888",
		DisableV2rayDeploy: false,
		Key:                "-",
		RemoveKey:          "-",
		LogLevel:           "error",
		LogPath:            defaultLogPath,
	}

	dc = &dbConfig{
		Database:   defaultDatabasePath,
		MaxPoolNum: defaultPoolNum,
	}

	if err := loadConfig(path.Join(rootPath, "config", "main.prod.config.yaml"), path.Join(rootPath, "config", "main.local.config.yaml"), &mc); nil != err {
		return err
	}

	if err := loadConfig(path.Join(rootPath, "config", "db.prod.config.yaml"), path.Join(rootPath, "config", "db.local.config.yaml"), &dc); nil != err {
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

func GetDbConfig() *dbConfig {
	if nil == dc {
		panic("获取了未初始化的配置")
	}

	return dc
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

func getRootPath() string {
	// 取到执行文件所在的目录作为根目录，否则在其他目录通过文件位置运行时会找不到配置文件
	executable, err := os.Executable()
	if nil != err {
		return ""
	}

	return filepath.Dir(executable)
}

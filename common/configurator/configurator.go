package configurator

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/util"
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

	mc = GetDefaultMailConfig()

	if err := mc.Load(filepath.Join(rootPath, "config", "main.config.yaml")); nil != err {
		return err
	}

	return nil
}

// GetDefaultMailConfig 获取默认的主配置文件
func GetDefaultMailConfig() *mainConfig {
	return &mainConfig{
		Address:          "0.0.0.0",
		Listen:           8888,
		EnableHttps:      false,
		GinReleaseMode:   true,
		Email:            fmt.Sprintf("%v@v2ray-helper.net", util.GenerateRandomString(16)),
		AllowV2rayDeploy: true,
		Key:              DefaultKey,
		RemoveKey:        DefaultRemoveKey,
		LogLevel:         "error",
	}
}

func GetMainConfig() *mainConfig {
	if nil == mc {
		panic("获取了未初始化的配置")
	}

	return mc
}

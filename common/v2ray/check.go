package v2ray

import (
	"errors"
	"fmt"
	"os"
)

const CmdPath = "/usr/local/bin/v2ray"

// CheckSystem 检查是否支持系统
func CheckSystem(goos, goArch string) bool {
	if "freebsd" == goos && ("386" == goArch || "amd64" == goArch) {
		return true
	}

	if "linux" == goos && ("386" == goArch || "amd64" == goArch || "arm" == goArch || "arm64" == goArch) {
		return true
	}

	if "darwin" == goos && ("amd64" == goArch || "arm64" == goArch) {
		return true
	}

	if "windows" == goos && ("386" == goArch || "amd64" == goArch || "arm" == goArch || "arm64" == goArch) {
		return true
	}

	return false
}

// CheckExists 检查系统内是否已存在V2ray
func CheckExists(v2rayPath string) (bool, error) {
	stat, err := os.Stat(v2rayPath)
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("无法获取文件信息: %v", err))
	}

	if stat.IsDir() {
		return false, nil
	}

	return true, nil
}

package v2ray

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

const v2rayCmdPath = "/usr/local/bin/v2ray"

// CheckSystem 检查是否支持系统
func CheckSystem() bool {
	if "freebsd" == runtime.GOOS && ("386" == runtime.GOARCH || "amd64" == runtime.GOARCH) {
		return true
	}

	if "linux" == runtime.GOOS && ("386" == runtime.GOARCH || "amd64" == runtime.GOARCH || "arm" == runtime.GOARCH || "arm64" == runtime.GOARCH) {
		return true
	}

	if "darwin" == runtime.GOOS && ("amd64" == runtime.GOARCH || "arm64" == runtime.GOARCH) {
		return true
	}

	if "windows" == runtime.GOOS && ("386" == runtime.GOARCH || "amd64" == runtime.GOARCH || "arm" == runtime.GOARCH || "arm64" == runtime.GOARCH) {
		return true
	}

	return false
}

// CheckExists 检查系统内是否已存在V2ray
func CheckExists() (bool, error) {
	stat, err := os.Stat(v2rayCmdPath)
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

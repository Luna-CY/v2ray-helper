package aria2

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const SystemdPath = "/etc/systemd/system/aria2.service"
const DefaultToken = "ARIARPCACCESSTOKEN"

const systemdConfig = `[Unit]
Description=Aria2 Service
Documentation=https://aria2.github.io/
After=network.target nss-lookup.target

[Service]
Type=simple
WorkingDirectory=/usr/local/v2ray-helper/temp/aria2
ExecStart=/usr/bin/aria2c --enable-rpc --rpc-secret=ARIARPCACCESSTOKEN
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target`

// Install 通过Apt或Yum安装Aria2
func Install(systemdPath, runtimePath string) error {
	uname, err := exec.Command("uname", "-a").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("检查系统类型失败: %v", err))
	}

	if strings.Contains(string(uname), "Debian") {
		res, err := exec.Command("sh", "-c", "apt update && apt install -y aria2").CombinedOutput()
		if nil != err && 0 < len(res) {
			return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
		}
	} else {
		res, err := exec.Command("sh", "-c", "yum update && yum install -y aria2").CombinedOutput()
		if nil != err && 0 < len(res) {
			return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
		}
	}

	// 安装到系统服务
	if err := os.MkdirAll(runtimePath, 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Dir(systemdPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
	}

	systemdFile, err := os.OpenFile(systemdPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
	}
	defer systemdFile.Close()

	if _, err := systemdFile.WriteString(systemdConfig); nil != err {
		return errors.New(fmt.Sprintf("安装Aria2失败: %v", err))
	}

	return nil
}

// InstallToSystem 安装到系统
func InstallToSystem() error {
	return Install(SystemdPath, filepath.Join(viper.GetString(configurator.KeyRootPath), "temp", "aria2"))
}

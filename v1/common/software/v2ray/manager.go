package v2ray

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning 检查是否在运行
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/bin/v2ray' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查V2ray运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Start 启动服务
func Start() error {
	_, err := exec.Command("service", "v2ray", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动V2ray服务失败: %v", err))
	}

	return nil
}

// Enable 设为开机启动
func Enable() error {
	_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/v2ray.service /etc/systemd/system/multi-user.target.wants/v2ray.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("V2ray服务设为开机启动失败: %v", err))
	}

	return nil
}

// Stop 停止服务
func Stop() error {
	_, err := exec.Command("service", "v2ray", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止V2ray服务失败: %v", err))
	}

	return nil
}

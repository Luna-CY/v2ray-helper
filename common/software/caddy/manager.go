package caddy

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning 检查Caddy是否在运行状态
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/bin/caddy' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Caddy运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Start 启动Caddy服务
func Start() error {
	_, err := exec.Command("service", "caddy", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动Caddy服务失败: %v", err))
	}

	return nil
}

// Enable 设置为开机启动
func Enable() error {
	_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/caddy.service /etc/systemd/system/multi-user.target.wants/caddy.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Caddy服务设为开机启动失败: %v", err))
	}

	return nil
}

// Stop 停止Caddy服务
func Stop() error {
	_, err := exec.Command("service", "caddy", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止Caddy服务失败: %v", err))
	}

	return nil
}

// Disable 取消开机自启
func Disable() error {
	_, err := exec.Command("sh", "-c", "rm -rf /etc/systemd/system/multi-user.target.wants/caddy.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Caddy服务取消开机启动失败: %v", err))
	}

	return nil
}

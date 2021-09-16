package aria2

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning 检查Aria2服务是否已启动
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/bin/aria2c' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Aria2运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Start 启动Aria2服务
func Start() error {
	_, err := exec.Command("service", "aria2", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动Aria2服务失败: %v", err))
	}

	return nil
}

// Stop 停止Aria2服务
func Stop() error {
	_, err := exec.Command("service", "aria2", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止Aria2服务失败: %v", err))
	}

	return nil
}

// Enable 设为开机自启
func Enable() error {
	_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/aria2.service /etc/systemd/system/multi-user.target.wants/aria2.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Aria2服务设为开机启动失败: %v", err))
	}

	return nil
}

// Disable 取消开机自启
func Disable() error {
	_, err := exec.Command("sh", "-c", "rm -rf /etc/systemd/system/multi-user.target.wants/aria2.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Aria2服务取消开机启动失败: %v", err))
	}

	return nil
}

package nginx

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// IsRunning 检查Nginx是否在运行状态
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/sbin/nginx' | grep -v grep | awk '{print $2}'").CombinedOutput()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Nginx状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Install Nginx通过Apt或Yum安装
func Install() error {
	uname, err := exec.Command("uname", "-a").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("检查系统类型失败: %v", err))
	}

	if strings.Contains(string(uname), "Debian") {
		res, err := exec.Command("sh", "-c", "apt update && apt install -y nginx").CombinedOutput()
		if nil != err && 0 < len(res) {
			return errors.New(fmt.Sprintf("安装Nginx失败: %v", err))
		}
	} else {
		res, err := exec.Command("sh", "-c", "yum update && yum install -y nginx").CombinedOutput()
		if nil != err && 0 < len(res) {
			return errors.New(fmt.Sprintf("安装Nginx失败: %v", err))
		}
	}

	return nil
}

// Start 启动Nginx服务
func Start() error {
	res, err := exec.Command("service", "nginx", "start").Output()
	if nil != err && 0 != len(res) {
		return errors.New(fmt.Sprintf("停止Nginx失败: %v", err))
	}

	return nil
}

// Stop 停止Nginx服务
func Stop() error {
	res, err := exec.Command("service", "nginx", "stop").Output()
	if nil != err && 0 != len(res) {
		return errors.New(fmt.Sprintf("停止Nginx失败: %v", err))
	}

	return nil
}

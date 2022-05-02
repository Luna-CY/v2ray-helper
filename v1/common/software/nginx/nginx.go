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

// Stop 停止Nginx服务
func Stop() error {
	res, err := exec.Command("service", "nginx", "stop").Output()
	if nil != err && 0 != len(res) {
		return errors.New(fmt.Sprintf("停止Nginx失败: %v", err))
	}

	return nil
}

// Disable 取消开机自启
func Disable() error {
	_, err := exec.Command("sh", "-c", "rm -rf /etc/systemd/system/multi-user.target.wants/nginx.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Nginx服务取消开机启动失败: %v", err))
	}

	return nil
}

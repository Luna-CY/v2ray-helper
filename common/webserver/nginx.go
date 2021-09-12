package webserver

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CheckNginxIsRunning 检查Nginx是否在运行状态
func CheckNginxIsRunning() (bool, error) {
	cmd := "ps ux | awk '/nginx/ && !/awk/ {print $2}'"

	res, err := exec.Command(cmd).Output()
	if nil != err {
		return false, errors.New(fmt.Sprintf("检查Nginx状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// StopNginx 停止Nginx进程
func StopNginx() error {
	cmd := "nginx -s stop"

	_, err := exec.Command(cmd).Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止Nginx失败: %v", err))
	}

	running, err := CheckNginxIsRunning()
	if nil != err {
		return err
	}

	if running {
		return errors.New("停止Nginx失败")
	}

	return nil
}

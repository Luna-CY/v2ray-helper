package vhelper

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/v2ray-helper/v2ray-helper' | grep -v grep | awk '{print $2}'").CombinedOutput()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查V2rayHelper状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

func Start() error {
	_, err := exec.Command("service", "v2ray-helper", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动V2rayHelper服务失败: %v", err))
	}

	return nil
}

func ReStart() error {
	_, err := exec.Command("service", "v2ray-helper", "restart").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("重启V2rayHelper服务失败: %v", err))
	}

	return nil
}

func Stop() error {
	_, err := exec.Command("service", "v2ray-helper", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止V2rayHelper服务失败: %v", err))
	}

	return nil
}

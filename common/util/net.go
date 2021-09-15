package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// GetPublicIpv4 获取外网IPv4地址
func GetPublicIpv4() (string, error) {
	res, err := http.Get("https://checkip.amazonaws.com")
	if nil != err {
		return "", errors.New(fmt.Sprintf("无法获取本机外网IP: %v", err))
	}
	defer res.Body.Close()

	ipBytes, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return "", errors.New(fmt.Sprintf("无法获取本机外网IP: %v", err))
	}

	return string(ipBytes), nil
}

func CheckLocalPortIsAllow(port int) (bool, error) {
	res, err := exec.Command("lsof", "-i", fmt.Sprintf(":%v", port)).Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查端口失败: %v", err))
	}

	if 0 < len(res) {
		return false, nil
	}

	return true, nil
}

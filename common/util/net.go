package util

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
)

// GetPublicIpv4 获取外网IPv4地址
func GetPublicIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("获取IP地址失败")

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

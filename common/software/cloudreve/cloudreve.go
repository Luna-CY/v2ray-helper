package cloudreve

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const Name = "cloudreve"
const InstallTo = "/usr/local/cloudreve"
const SystemdPath = "/etc/systemd/system/cloudreve.service"

const lastReleaseApi = "https://api.github.com/repos/cloudreve/Cloudreve/releases/latest"
const downloadUrlTemplate = "https://github.com/cloudreve/Cloudreve/releases/download/%v/cloudreve_%v_%v_%v.tar.gz"
const systemdConfig = `[Unit]
Description=Cloudreve
Documentation=https://docs.cloudreve.org
After=network.target nss-lookup.target

[Service]
Type=simple
WorkingDirectory=/usr/local/cloudreve
ExecStart=/usr/local/cloudreve/cloudreve
Restart=on-failure
RestartPreventExitStatus=23
RestartSec=5s

[Install]
WantedBy=multi-user.target`

// IsRunning 检查Cloudreve是否在运行状态
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/cloudreve/cloudreve' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Cloudreve运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Install 安装Cloudreve服务
func Install(goos, goArch, version, installTo, systemdPath string) error {
	downloadUrl, err := GetDownloadUrl(goos, goArch, version)
	if nil != err {
		return err
	}

	if err := os.MkdirAll(installTo, 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}

	res, err := http.Get(downloadUrl)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}
	defer res.Body.Close()

	gr, err := gzip.NewReader(res.Body)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}

	// 安装二进制可执行文件
	tr := tar.NewReader(gr)
	for {
		f, err := tr.Next()
		if nil != err {
			if io.EOF == err {
				break
			} else {
				return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
			}
		}

		if f.FileInfo().IsDir() {
			continue
		}

		var tf *os.File

		if "cloudreve" == f.Name {
			tf, err = os.OpenFile(filepath.Join(installTo, f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if nil != err {
				return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
			}
		} else {
			continue
		}

		if _, err := io.Copy(tf, tr); nil != err {
			return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
		}

		tf.Close()
	}

	// 安装到系统服务
	if err := os.MkdirAll(filepath.Dir(systemdPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}

	systemdFile, err := os.OpenFile(systemdPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}
	defer systemdFile.Close()

	if _, err := systemdFile.WriteString(systemdConfig); nil != err {
		return errors.New(fmt.Sprintf("安装Cloudreve失败: %v", err))
	}

	return nil
}

// InstallLastRelease 安装最新版本
func InstallLastRelease() error {
	version, err := GetLastReleaseVersion()
	if nil != err {
		return err
	}

	return Install(runtime.GOOS, runtime.GOARCH, version, InstallTo, SystemdPath)
}

// GetDownloadUrl 获取版本的下载链接
func GetDownloadUrl(goos, goArch, version string) (string, error) {
	if "linux" != goos {
		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "amd64" == goArch {
		return fmt.Sprintf(downloadUrlTemplate, version, version, goos, "amd64"), nil
	}

	if "arm" == goArch {
		return fmt.Sprintf(downloadUrlTemplate, version, version, goos, "arm"), nil
	}

	if "arm64" == goArch {
		return fmt.Sprintf(downloadUrlTemplate, version, version, goos, "arm64"), nil
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
}

// GetLastReleaseVersion 获取Cloudreve最后一个版本的版本号
func GetLastReleaseVersion() (string, error) {
	res, err := http.Get(lastReleaseApi)
	if nil != err {
		return "", errors.New(fmt.Sprintf("查询版本数据失败: %v", err))
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return "", errors.New(fmt.Sprintf("无法读取版本数据: %v", err))
	}

	var response lastReleaseVersionStruct
	if err := json.Unmarshal(content, &response); nil != err {
		return "", errors.New(fmt.Sprintf("无法解析版本数据: %v", err))
	}

	return strings.Trim(response.TagName, "v"), nil
}

// Start 启动Cloudreve服务
func Start() error {
	_, err := exec.Command("service", "cloudreve", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动Cloudreve服务失败: %v", err))
	}

	return nil
}

// Enable 设置为开机启动
func Enable() error {
	_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/cloudreve.service /etc/systemd/system/multi-user.target.wants/cloudreve.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Cloudreve服务设为开机启动失败: %v", err))
	}

	return nil
}

// Stop 停止Cloudreve服务
func Stop() error {
	_, err := exec.Command("service", "cloudreve", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止Cloudreve服务失败: %v", err))
	}

	return nil
}

// Disable 取消开机自启
func Disable() error {
	_, err := exec.Command("sh", "-c", "rm -rf /etc/systemd/system/multi-user.target.wants/cloudreve.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Cloudreve服务取消开机启动失败: %v", err))
	}

	return nil
}

// ResetAdminPassword 重置并获取管理员密码
func ResetAdminPassword() (string, error) {
	res, err := exec.Command("sh", "-c", "/usr/local/cloudreve/cloudreve --database-script ResetAdminPassword").Output()
	if nil != err {
		return "", errors.New(fmt.Sprintf("重置Cloudreve管理员密码失败: %v", err))
	}

	reg := regexp.MustCompile("^.*初始管理员密码已更改为：.{8}$")

	output := strings.Split(string(res), "\n")
	for _, line := range output {
		if !strings.HasPrefix(line, "[Info]") {
			continue
		}

		if reg.MatchString(line) {
			return line[len(line)-8:], nil
		}
	}

	return "", errors.New(fmt.Sprintf("获取管理员密码失败"))
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

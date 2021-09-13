package webserver

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
	"runtime"
	"strings"
)

const CaddyInstallTo = "/usr/local/bin"
const CaddyConfigPath = "/etc/caddy/Caddyfile"
const CaddySystemdPath = "/etc/systemd/system/caddy.service"

const caddyLastReleaseApi = "https://api.github.com/repos/caddyserver/caddy/releases/latest"
const caddyDownloadUrlTemplate = "https://github.com/caddyserver/caddy/releases/download/v%v/caddy_%v_%v_%v.tar.gz"

const caddyConfigTemplate = `%v:%v {
    reverse_proxy %v 127.0.0.1:%v
}`
const caddyAndCloudConfigTemplate = `%v:%v {
    reverse_proxy %v 127.0.0.1:%v

    reverse_proxy / 127.0.0.1:5212
}`
const caddySystemdConfig = `[Unit]
Description=Caddy Service
Documentation=https://caddyserver.com/docs/
After=network.target nss-lookup.target

[Service]
Type=simple
ExecStart=/usr/local/bin/caddy run -config /etc/caddy/Caddyfile
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target`

// CaddyIsInstall 检查是否已安装Caddy
func CaddyIsInstall() (bool, error) {
	stat, err := os.Stat(filepath.Join(CaddyInstallTo, "caddy"))
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("检查Caddy失败: %v", err))
	}

	if stat.IsDir() {
		return false, nil
	}

	return true, nil
}

// InstallCaddy 安装Caddy
func InstallCaddy(goos, goArch, version, installTo, configPath, systemdPath string) error {
	downloadUrl, err := GetCaddyDownloadUrl(goos, goArch, version)
	if nil != err {
		return err
	}

	if err := os.MkdirAll(installTo, 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	res, err := http.Get(downloadUrl)
	if nil != err {
		return errors.New(fmt.Sprintf("下载Caddy失败: %v", err))
	}
	defer res.Body.Close()

	gr, err := gzip.NewReader(res.Body)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	// 安装二进制可执行文件
	tr := tar.NewReader(gr)
	for {
		f, err := tr.Next()
		if nil != err {
			if io.EOF == err {
				break
			} else {
				return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
			}
		}

		if f.FileInfo().IsDir() {
			continue
		}

		var tf *os.File

		if "caddy" == f.Name {
			tf, err = os.OpenFile(filepath.Join(installTo, f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if nil != err {
				return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
			}
		} else {
			continue
		}

		if _, err := io.Copy(tf, tr); nil != err {
			return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
		}

		tf.Close()
	}

	// 安装到系统服务
	if err := os.MkdirAll(filepath.Dir(systemdPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	systemdFile, err := os.OpenFile(systemdPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}
	defer systemdFile.Close()

	if _, err := systemdFile.WriteString(caddySystemdConfig); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	// 重置配置文件
	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}
	defer configFile.Close()

	return nil
}

// InstallCaddyLastRelease 安装最新版本
func InstallCaddyLastRelease() error {
	version, err := GetCaddyLastReleaseVersion()
	if nil != err {
		return err
	}

	return InstallCaddy(runtime.GOOS, runtime.GOARCH, version, CaddyInstallTo, CaddyConfigPath, CaddySystemdPath)
}

// AppendCaddyConfig 添加Caddy的配置
func AppendCaddyConfig(configPath, host string, useTls bool, CaddyPort int, cloud bool, path string) error {
	port := 443
	if !useTls {
		port = 80
	}

	var content string
	if cloud {
		if "/" == path {
			return errors.New(fmt.Sprintf("参数错误，同时配置cloudreve时，path参数不能为/"))
		}

		content = fmt.Sprintf(caddyAndCloudConfigTemplate, host, port, path, CaddyPort)
	} else {
		content = fmt.Sprintf(caddyConfigTemplate, host, port, path, CaddyPort)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("配置Caddy失败: %v", err))
	}

	cf, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
	}
	defer cf.Close()

	if _, err := cf.WriteString(content); nil != err {
		return errors.New(fmt.Sprintf("无法写入配置文件: %v", err))
	}

	return nil
}

// AppendCaddyConfigOnlyV2rayToSystem 只添加Caddy的配置到系统
func AppendCaddyConfigOnlyV2rayToSystem(host string, useTls bool, CaddyPort int, path string) error {
	return AppendCaddyConfig(CaddyConfigPath, host, useTls, CaddyPort, false, path)
}

// AppendCaddyConfigCaddyAndCloudToSystem 添加Caddy与Cloudreve的配置到系统
func AppendCaddyConfigCaddyAndCloudToSystem(host string, useTls bool, CaddyPort int, path string) error {
	return AppendCaddyConfig(CaddyConfigPath, host, useTls, CaddyPort, true, path)
}

// GetCaddyLastReleaseVersion 获取Caddy最后一个版本的版本号
func GetCaddyLastReleaseVersion() (string, error) {
	res, err := http.Get(caddyLastReleaseApi)
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

// GetCaddyDownloadUrl 获取版本的下载链接
func GetCaddyDownloadUrl(goos, goArch, version string) (string, error) {
	if "linux" != goos {
		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "amd64" == goArch {
		return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "amd64"), nil
	}

	if "arm64" == goArch {
		return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "arm64"), nil
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
}

// CheckCaddyIsRunning 检查Caddy是否在运行状态
func CheckCaddyIsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/bin/caddy' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Caddy运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// StartCaddy 启动Caddy服务
func StartCaddy() error {
	_, err := exec.Command("service", "caddy", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动Caddy服务失败: %v", err))
	}

	return nil
}

// ReStartCaddy 重启Caddy服务
func ReStartCaddy() error {
	_, err := exec.Command("service", "caddy", "restart").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("重启Caddy服务失败: %v", err))
	}

	return nil
}

// StopCaddy 停止Caddy服务
func StopCaddy() error {
	_, err := exec.Command("service", "caddy", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("重启Caddy服务失败: %v", err))
	}

	return nil
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

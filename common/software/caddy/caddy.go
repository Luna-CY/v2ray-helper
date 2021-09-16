package caddy

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
	"path/filepath"
	"runtime"
	"strings"
)

const InstallTo = "/usr/local/bin"
const ConfigPath = "/etc/caddy/Caddyfile"
const SystemdPath = "/etc/systemd/system/caddy.service"

const lastReleaseApi = "https://api.github.com/repos/caddyserver/caddy/releases/latest"
const downloadUrlTemplate = "https://github.com/caddyserver/caddy/releases/download/v%v/caddy_%v_%v_%v.tar.gz"

const systemdConfig = `[Unit]
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

// Install 安装Caddy
func Install(goos, goArch, version, installTo, configPath, systemdPath string) error {
	downloadUrl, err := GetDownloadUrl(goos, goArch, version)
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

	if _, err := systemdFile.WriteString(systemdConfig); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	// 添加默认配置文件
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Join(filepath.Dir(configPath), "sites-enabled"), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	if err := os.RemoveAll(configPath); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}
	defer configFile.Close()

	if _, err := configFile.WriteString("import sites-enabled/*"); nil != err {
		return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
	}

	return nil
}

// InstallLastRelease 安装最新版本
func InstallLastRelease() error {
	version, err := GetLastReleaseVersion()
	if nil != err {
		return err
	}

	return Install(runtime.GOOS, runtime.GOARCH, version, InstallTo, ConfigPath, SystemdPath)
}

// GetLastReleaseVersion 获取Caddy最后一个版本的版本号
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

// GetDownloadUrl 获取版本的下载链接
func GetDownloadUrl(goos, goArch, version string) (string, error) {
	if "linux" != goos {
		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "amd64" == goArch {
		return fmt.Sprintf(downloadUrlTemplate, version, version, goos, "amd64"), nil
	}

	if "arm64" == goArch {
		return fmt.Sprintf(downloadUrlTemplate, version, version, goos, "arm64"), nil
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

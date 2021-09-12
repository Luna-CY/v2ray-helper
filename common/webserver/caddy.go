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
	"path/filepath"
	"runtime"
	"strings"
)

const CaddyInstallTo = "/usr/local/bin"
const CaddyConfigPath = "/etc/caddy/Caddyfile"

const caddyLastReleaseApi = "https://api.github.com/repos/caddyserver/caddy/releases/latest"
const caddyDownloadUrlTemplate = "https://github.com/caddyserver/caddy/releases/download/v%v/caddy_%v_%v_%v.tar.gz"

const caddyConfigTemplate = `\n%v:%v {
    reverse_proxy %v 127.0.0.1:%v
}`
const caddyAndCloudConfigTemplate = `\n%v:%v {
    reverse_proxy %v 127.0.0.1:%v

    reverse_proxy / 127.0.0.1:5212
}`

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
func InstallCaddy(goos, goArch, version, installTo string) error {
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
		fmt.Println(f.Name)

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

	return nil
}

// InstallCaddyLastRelease 安装最新版本
func InstallCaddyLastRelease() error {
	version, err := GetCaddyLastReleaseVersion()
	if nil != err {
		return err
	}

	return InstallCaddy(runtime.GOOS, runtime.GOARCH, version, CaddyInstallTo)
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
	if "freebsd" == goos {
		if "amd64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "amd64"), nil
		}

		if "arm64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "arm64"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "linux" == goos {
		if "amd64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "amd64"), nil
		}

		if "arm64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, goos, "arm64"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "darwin" == goos {
		if "amd64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, "mac", "amd64"), nil
		}

		if "arm64" == goArch {
			return fmt.Sprintf(caddyDownloadUrlTemplate, version, version, "mac", "arm64"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

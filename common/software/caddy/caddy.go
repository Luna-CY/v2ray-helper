package caddy

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	vhr "github.com/Luna-CY/v2ray-helper/common/runtime"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const InstallTo = "/usr/local/bin"
const ConfigPath = "/etc/caddy/Caddyfile"
const SystemdPath = "/etc/systemd/system/caddy.service"

const lastReleaseApi = "https://api.github.com/repos/caddyserver/caddy/releases/latest"
const downloadUrlTemplate = "https://github.com/caddyserver/caddy/releases/download/v%v/caddy_%v_%v_%v.tar.gz"

const configTemplate = `%v:%v {
    reverse_proxy %v 127.0.0.1:%v
}`
const configTlsTemplate = `%v:%v {
    tls %v %v

    reverse_proxy %v 127.0.0.1:%v {
        transport http {
            versions h2c
        }
    }
}`
const configCloudreveTemplate = `%v:%v {
    reverse_proxy %v 127.0.0.1:%v

    reverse_proxy 127.0.0.1:5212
}`
const configCloudreveAndTlsTemplate = `%v:%v {
    tls %v %v

    reverse_proxy %v 127.0.0.1:%v {
        transport http {
            versions h2c
        }
    }

    reverse_proxy 127.0.0.1:5212
}`
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

// IsInstall 检查是否已安装Caddy
func IsInstall() (bool, error) {
	stat, err := os.Stat(filepath.Join(InstallTo, "caddy"))
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

	// 如果配置文件不存在就创建一个新的配置文件
	stat, err := os.Stat(configPath)
	if nil != err {
		if !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
		}

		configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if nil != err {
			return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
		}
		defer configFile.Close()
	} else {
		if stat.IsDir() {
			if err := os.RemoveAll(configPath); nil != err {
				return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
			}
		}

		configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if nil != err {
			return errors.New(fmt.Sprintf("安装Caddy失败: %v", err))
		}
		defer configFile.Close()
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

// AppendConfig 添加Caddy的配置
func AppendConfig(configPath, host string, useTls bool, v2rayPort int, path string, cloudreve, cusTls bool) error {
	port := 443
	if !useTls {
		port = 80
	}

	content := fmt.Sprintf(configTemplate, host, port, path, v2rayPort)

	// 自定义证书
	if cusTls {
		key := filepath.Join(vhr.GetRootPath(), certificate.CertDirName, host, "private.key")
		cert := filepath.Join(vhr.GetRootPath(), certificate.CertDirName, host, "cert.pem")

		content = fmt.Sprintf(configTlsTemplate, host, port, cert, key, path, v2rayPort)
	}

	// Cloudreve
	if cloudreve {
		if "/" == path {
			return errors.New(fmt.Sprintf("参数错误，同时配置cloudreve时，path参数不能为/"))
		}

		content = fmt.Sprintf(configCloudreveTemplate, host, port, path, v2rayPort)
	}

	// Cloudreve与自定义证书
	if cloudreve && cusTls {
		if "/" == path {
			return errors.New(fmt.Sprintf("参数错误，同时配置cloudreve时，path参数不能为/"))
		}

		key := filepath.Join(vhr.GetRootPath(), certificate.CertDirName, host, "private.key")
		cert := filepath.Join(vhr.GetRootPath(), certificate.CertDirName, host, "cert.pem")

		content = fmt.Sprintf(configCloudreveAndTlsTemplate, host, port, cert, key, path, v2rayPort)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("配置Caddy失败: %v", err))
	}

	cf, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
	}
	defer cf.Close()

	if _, err := cf.WriteString(content); nil != err {
		return errors.New(fmt.Sprintf("无法写入配置文件: %v", err))
	}

	return nil
}

// AppendConfigOnlyV2rayToSystem 只添加V2ray的配置到系统
func AppendConfigOnlyV2rayToSystem(host string, useTls bool, v2rayPort int, path string) error {
	return AppendConfig(ConfigPath, host, useTls, v2rayPort, path, false, false)
}

// AppendConfigV2rayAndHTTP2ToSystem 添加V2ray的配置与HTTP2到系统
func AppendConfigV2rayAndHTTP2ToSystem(host string, useTls bool, v2rayPort int, path string) error {
	return AppendConfig(ConfigPath, host, useTls, v2rayPort, path, false, true)
}

// AppendConfigV2rayAndCloudreveToSystem 添加V2ray与Cloudreve的配置到系统
func AppendConfigV2rayAndCloudreveToSystem(host string, useTls bool, v2rayPort int, path string) error {
	return AppendConfig(ConfigPath, host, useTls, v2rayPort, path, true, false)
}

// AppendConfigV2rayAndHTTP2AndCloudreveToSystem 添加V2ray的HTTP2模式与Cloudreve两个配置到系统
func AppendConfigV2rayAndHTTP2AndCloudreveToSystem(host string, useTls bool, v2rayPort int, path string) error {
	return AppendConfig(ConfigPath, host, useTls, v2rayPort, path, true, true)
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

// IsRunning 检查Caddy是否在运行状态
func IsRunning() (bool, error) {
	res, err := exec.Command("sh", "-c", "ps -ef | grep '/usr/local/bin/caddy' | grep -v grep | awk '{print $2}'").Output()
	if nil != err && 0 != len(res) {
		return false, errors.New(fmt.Sprintf("检查Caddy运行状态失败: %v", err))
	}

	return "" != strings.TrimSpace(string(res)), nil
}

// Start 启动Caddy服务
func Start() error {
	_, err := exec.Command("service", "caddy", "start").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("启动Caddy服务失败: %v", err))
	}

	return nil
}

// Enable 设置为开机启动
func Enable() error {
	_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/caddy.service /etc/systemd/system/multi-user.target.wants/caddy.service").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("Caddy服务设为开机启动失败: %v", err))
	}

	return nil
}

// Stop 停止Caddy服务
func Stop() error {
	_, err := exec.Command("service", "caddy", "stop").Output()
	if nil != err {
		return errors.New(fmt.Sprintf("停止Caddy服务失败: %v", err))
	}

	return nil
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

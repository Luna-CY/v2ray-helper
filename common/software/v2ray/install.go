package v2ray

import (
	"archive/zip"
	"bytes"
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
const SharePath = "/usr/local/share"
const SystemdPath = "/etc"

// Install 安装指定版本
func Install(goos, goArch, version, installTo, sharePath, systemdPath string) error {
	downloadUrl, err := GetDownloadUrl(goos, goArch, version)
	if nil != err {
		return err
	}

	if err := os.MkdirAll(installTo, 0755); nil != err {
		return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Join(sharePath, "v2ray"), 0755); nil != err {
		return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
	}

	res, err := http.Get(downloadUrl)
	if nil != err {
		return errors.New(fmt.Sprintf("下载V2ray失败: %v", err))
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
	}

	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if nil != err {
		return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
	}

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		zf, err := f.Open()
		if nil != err {
			return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
		}

		var tf *os.File

		if "v2ray" == f.Name || "v2ctl" == f.Name {
			tf, err = os.OpenFile(filepath.Join(installTo, f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}
		} else if "geoip.dat" == f.Name || "geosite.dat" == f.Name {
			tf, err = os.OpenFile(filepath.Join(sharePath, "v2ray", f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}
		} else if strings.HasPrefix(f.Name, "systemd/system") {
			dstPath := filepath.Join(systemdPath, f.Name)
			if err := os.MkdirAll(filepath.Dir(dstPath), 0755); nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}

			tf, err = os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}
		} else {
			continue
		}

		if _, err := io.Copy(tf, zf); nil != err {
			return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
		}

		zf.Close()
		tf.Close()
	}

	return nil
}

// InstallLastRelease 安装最新版本
func InstallLastRelease() error {
	version, err := GetLastReleaseVersion()
	if nil != err {
		return err
	}

	return Install(runtime.GOOS, runtime.GOARCH, version, InstallTo, SharePath, SystemdPath)
}

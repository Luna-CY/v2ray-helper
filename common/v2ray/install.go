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
	"path"
	"strings"
)

const installDir = "/usr/local/bin"
const datPath = "/usr/local/share/v2ray"

// Install 安装指定版本
func Install(version string) error {
	downloadUrl, err := GetDownloadUrl(version)
	if nil != err {
		return err
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

	if err := os.MkdirAll(path.Join(datPath, "v2ray"), 0755); nil != err {
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
			tf, err = os.OpenFile(path.Join(installDir, f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}
		} else if "geoip.dat" == f.Name || "geosite.dat" == f.Name {
			tf, err = os.OpenFile(path.Join(datPath, "v2ray", f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if nil != err {
				return errors.New(fmt.Sprintf("安装V2ray失败: %v", err))
			}
		} else if strings.HasPrefix(f.Name, "systemd/system") {
			tf, err = os.OpenFile(path.Join("/etc", f.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

	return Install(version)
}

package v2ray

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const lastReleaseApi = "https://api.github.com/repos/v2fly/v2ray-core/releases/latest"
const downloadUrlTemplate = "https://github.com/v2fly/v2ray-core/releases/download/v%v/v2ray-%v-%v.zip"

// GetLastReleaseVersion 获取V2ray最后一个版本的版本号
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
	if "freebsd" == goos {
		if "386" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "32"), nil
		}

		if "amd64" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "64"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "linux" == goos {
		if "386" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "32"), nil
		}

		if "amd64" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "64"), nil
		}

		if "arm32" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "arm32-v8a"), nil
		}

		if "arm64" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, goos, "arm64-v8a"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	if "darwin" == goos {
		if "amd64" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, "macos", "64"), nil
		}

		if "arm64" == goArch {
			return fmt.Sprintf(downloadUrlTemplate, version, "macos", "arm64-v8a"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", goos, goArch))
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

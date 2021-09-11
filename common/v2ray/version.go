package v2ray

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
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
func GetDownloadUrl(version string) (string, error) {
	if "freebsd" == runtime.GOOS {
		if "386" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "32"), nil
		}

		if "amd64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "64"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH))
	}

	if "linux" == runtime.GOOS {
		if "386" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "32"), nil
		}

		if "amd64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "64"), nil
		}

		if "arm32" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "arm32-v8a"), nil
		}

		if "arm64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "arm64-v8a"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH))
	}

	if "windows" == runtime.GOOS {
		if "386" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "32"), nil
		}

		if "amd64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "64"), nil
		}

		if "arm32" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "arm32-v7a"), nil
		}

		if "arm64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, runtime.GOOS, "arm64-v8a"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH))
	}

	if "darwin" == runtime.GOOS {
		if "amd64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, "macos", "64"), nil
		}

		if "arm64" == runtime.GOARCH {
			return fmt.Sprintf(downloadUrlTemplate, version, "macos", "arm64-v8a"), nil
		}

		return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH))
	}

	return "", errors.New(fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH))
}

type lastReleaseVersionStruct struct {
	TagName string `json:"tag_name"`
}

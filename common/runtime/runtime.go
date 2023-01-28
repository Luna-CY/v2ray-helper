package runtime

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"os"
	"path/filepath"
)

// AbsRootPath 获取项目运行绝对目录
func AbsRootPath(homeDir string) string {
	if "" != homeDir {
		if filepath.IsAbs(homeDir) {
			return homeDir
		}

		if "." != homeDir {
			absPath, err := filepath.Abs(homeDir)
			if nil == err {
				return absPath
			}
		}
	}

	// 取到执行文件所在的目录作为根目录，否则在其他目录通过文件位置运行时会找不到配置文件
	executable, err := os.Executable()
	if nil != err {
		return ""
	}

	return filepath.Dir(executable)
}

const AcmeCertDirName = "acme"

// GetAcmeCertificatePath 获取证书目录的路径
func GetAcmeCertificatePath() string {
	return filepath.Join(configurator.Configure.Home, AcmeCertDirName)
}

const V2rayConfig = "v2ray.json"

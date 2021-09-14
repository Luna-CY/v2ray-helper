package certificate

import (
	"crypto"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/go-acme/lego/v4/registration"
	"os"
	"path/filepath"
)

const CertDirName = "certs"

// Init 初始化证书管理器
func Init(rootPath string) error {
	if err := os.MkdirAll(filepath.Join(rootPath, CertDirName), 0755); nil != err {
		return errors.New(fmt.Sprintf("初始化证书环境失败: %v", err))
	}

	return nil
}

// CertIsExists 检查证书是否存在
func CertIsExists(host string) (bool, error) {
	path := filepath.Join(runtime.GetRootPath(), CertDirName, host)
	hostStat, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if !hostStat.IsDir() {
		return false, errors.New(fmt.Sprintf("证书配置错误，请检查目录: %v", path))
	}

	keyStat, err := os.Stat(filepath.Join(path, "private.key"))
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if keyStat.IsDir() {
		return false, errors.New(fmt.Sprintf("证书配置错误，请检查目录: %v", filepath.Join(path, "private.key")))
	}

	certStat, err := os.Stat(filepath.Join(path, "cert.pem"))
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if certStat.IsDir() {
		return false, errors.New(fmt.Sprintf("证书配置错误，请检查目录: %v", filepath.Join(path, "cert.pem")))
	}

	return true, nil
}

type user struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *user) GetEmail() string {
	return u.Email
}

func (u *user) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *user) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

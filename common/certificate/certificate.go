package certificate

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/go-acme/lego/v4/registration"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Init 初始化证书管理器
func Init(ctx context.Context) error {
	if err := os.MkdirAll(runtime.GetCertificatePath(), 0755); nil != err {
		return errors.New(fmt.Sprintf("初始化证书环境失败: %v", err))
	}

	return initManager(ctx)
}

func newUser(email string) (*user, error) {
	u := &user{Email: email}

	privateKeyPath := filepath.Join(viper.GetString(configurator.KeyRootPath), fmt.Sprintf("%v.key", email))
	k, err := os.Open(privateKeyPath)
	if nil != err {
		if !os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("无法打开用户的私钥文件: %v", err))
		}
	}

	if nil != k {
		defer k.Close()

		content, err := io.ReadAll(k)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("无法读取用户的私钥文件: %v", err))
		}

		pp, _ := pem.Decode(content)
		pk, err := x509.ParseECPrivateKey(pp.Bytes)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("无法解析用户的私钥内容: %v", err))
		}

		u.key = pk
	}

	if nil == u.key {
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("无法生成私钥: %v", err))
		}

		f, err := os.Create(privateKeyPath)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("无法创建用户的私钥文件: %v", err))
		}
		defer f.Close()

		pp, err := x509.MarshalECPrivateKey(pk)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("序列化用户私钥失败: %v", err))
		}

		if err := pem.Encode(f, &pem.Block{Type: "PRIVATE KEY", Bytes: pp}); nil != err {
			return nil, errors.New(fmt.Sprintf("编码用户私钥内容失败: %v", err))
		}

		u.key = pk
	}

	return u, nil
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

type Certificate struct {
	host string

	privateKeyContent       []byte
	certificateContent      []byte
	csrContent              []byte
	issueCertificateContent []byte

	certificate *x509.Certificate
}

func (c *Certificate) GetPrivateKeyContent() []byte {
	return c.privateKeyContent
}

func (c *Certificate) GetPrivateKeyFilePath() string {
	return filepath.Join(runtime.GetCertificatePath(), c.host, "private.key")
}

func (c *Certificate) GetCertificateContent() []byte {
	return c.certificateContent
}

func (c *Certificate) GetCertificateFilePath() string {
	return filepath.Join(runtime.GetCertificatePath(), c.host, "cert.pem")
}

func (c *Certificate) GetCsrContent() []byte {
	return c.csrContent
}

func (c *Certificate) GetCsrFilePath() string {
	return filepath.Join(runtime.GetCertificatePath(), c.host, "cert.csr")
}

func (c *Certificate) GetIssueCertificate() []byte {
	return c.issueCertificateContent
}

func (c *Certificate) GetIssueCertificateFilePath() string {
	return filepath.Join(runtime.GetCertificatePath(), c.host, "cert.issue")
}

func (c *Certificate) GetExpireTime() time.Time {
	return c.certificate.NotAfter
}

// newCertificate 新建一个证书结构并解析证书
func newCertificate(host string) (*Certificate, error) {
	path := filepath.Join(runtime.GetCertificatePath())

	privateKeyFile, err := os.Open(filepath.Join(path, host, "private.key"))
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法打开私钥文件: %v", err))
	}

	privateKeyContent, err := ioutil.ReadAll(privateKeyFile)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法读取私钥文件: %v", err))
	}
	defer privateKeyFile.Close()

	certFile, err := os.Open(filepath.Join(path, host, "cert.pem"))
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法打开证书文件: %v", err))
	}

	certContent, err := ioutil.ReadAll(certFile)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法读取证书文件: %v", err))
	}
	defer certFile.Close()

	csrFile, err := os.Open(filepath.Join(path, host, "cert.csr"))
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法打开证书文件: %v", err))
	}

	csrContent, err := ioutil.ReadAll(csrFile)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法读取证书文件: %v", err))
	}
	defer csrFile.Close()

	issueFile, err := os.Open(filepath.Join(path, host, "cert.issue"))
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法打开证书文件: %v", err))
	}

	issueContent, err := ioutil.ReadAll(issueFile)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法读取证书文件: %v", err))
	}
	defer issueFile.Close()

	cert := new(Certificate)
	cert.host = host
	cert.privateKeyContent = privateKeyContent
	cert.certificateContent = certContent
	cert.csrContent = csrContent
	cert.issueCertificateContent = issueContent

	block, _ := pem.Decode(certContent)
	certificate, err := x509.ParseCertificate(block.Bytes)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("解析证书失败: %v", err))
	}

	cert.certificate = certificate

	return cert, nil
}

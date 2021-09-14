package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
	"path/filepath"
)

// IssueNew 申请新的证书
// 返回的第一个参数是证书的私钥内容
// 返回的第二个参数是证书的内容
func IssueNew(host, email string) ([]byte, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("无法生成私钥: %v", err))
	}

	userEntry := &user{Email: email, key: privateKey}
	config := lego.NewConfig(userEntry)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("无法创建Acme客户端: %v", err))
	}

	if err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80")); err != nil {
		return nil, nil, errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	if err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443")); err != nil {
		return nil, nil, errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("注册用户失败: %v", err))
	}

	userEntry.Registration = reg

	request := certificate.ObtainRequest{Domains: []string{host}, Bundle: true}
	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	if err := os.RemoveAll(filepath.Join(runtime.GetRootPath(), CertDirName, host)); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Join(runtime.GetRootPath(), CertDirName, host), 0755); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存私钥
	keyFile, err := os.OpenFile(filepath.Join(runtime.GetRootPath(), CertDirName, host, "private.key"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer keyFile.Close()

	if _, err := keyFile.Write(cert.PrivateKey); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存证书
	certFile, err := os.OpenFile(filepath.Join(runtime.GetRootPath(), CertDirName, host, "cert.pem"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer certFile.Close()

	if _, err := certFile.Write(cert.Certificate); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存csr
	csrFile, err := os.OpenFile(filepath.Join(runtime.GetRootPath(), CertDirName, host, "cert.csr"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer csrFile.Close()

	if _, err := csrFile.Write(cert.CSR); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存IssuerCertificate
	issueFile, err := os.OpenFile(filepath.Join(runtime.GetRootPath(), CertDirName, host, "cert.issue"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer issueFile.Close()

	if _, err := issueFile.Write(cert.IssuerCertificate); nil != err {
		return nil, nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	return cert.PrivateKey, cert.Certificate, nil
}

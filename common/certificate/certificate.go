package certificate

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type user struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *user) GetEmail() string {
	return u.Email
}

func (u *user) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *user) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

// IssueNew 申请新的证书
// 返回的第一个参数是证书的私钥内容
// 返回的第二个参数是证书的内容
func IssueNew(host, email string) ([]byte, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("无法生成私钥: %v", err))
	}

	userEntry := &user{Email: email, Key: privateKey}
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

	return cert.PrivateKey, cert.Certificate, nil
}

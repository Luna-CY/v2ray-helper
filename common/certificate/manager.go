package certificate

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/Luna-CY/v2ray-helper/common/software/caddy"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
	"github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var manager *Manager

// initManager 初始化证书管理器，读取所有证书
func initManager(ctx context.Context) error {
	cm := new(Manager)
	cm.ctx = ctx
	cm.certs = map[string]*Certificate{}

	path := runtime.GetCertificatePath()
	stat, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("证书目录不存在"))
		}

		return errors.New(fmt.Sprintf("初始化证书管理器失败: %v", err))
	}

	if !stat.IsDir() {
		return errors.New(fmt.Sprintf("初始化证书管理器失败: 证书路径不是一个目录: %v", path))
	}

	hostList, err := ioutil.ReadDir(path)
	if nil != err {
		return errors.New(fmt.Sprintf("初始化证书管理器失败: %v", err))
	}

	cm.Lock()
	defer cm.Unlock()

	for _, host := range hostList {
		if !host.IsDir() {
			return errors.New(fmt.Sprintf("检查证书失败: 此路径不是一个目录: %v", filepath.Join(path, host.Name())))
		}

		cert, err := loadCertificate(host.Name())
		if nil != err {
			return err
		}

		cm.certs[host.Name()] = cert
	}

	manager = cm

	return nil
}

// loadCertificate 加载证书
func loadCertificate(host string) (*Certificate, error) {
	path := runtime.GetCertificatePath()

	keyStat, err := os.Stat(filepath.Join(path, host, "private.key"))
	if nil != err {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("检查证书失败: 无法找到私钥文件: %v", filepath.Join(path, host, "private.key")))
		}

		return nil, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if keyStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("证书配置错误，此路径不是一个有效的文件: %v", filepath.Join(path, host, "private.key")))
	}

	certStat, err := os.Stat(filepath.Join(path, host, "cert.pem"))
	if nil != err {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("检查证书失败: 无法找到证书文件: %v", filepath.Join(path, host, "cert.pem")))
		}

		return nil, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if certStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("证书配置错误，此路径不是一个有效的文件: %v", filepath.Join(path, host, "cert.pem")))
	}

	csrStat, err := os.Stat(filepath.Join(path, host, "cert.csr"))
	if nil != err {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("检查证书失败: 无法找到证书文件: %v", filepath.Join(path, host, "cert.csr")))
		}

		return nil, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if csrStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("证书配置错误，此路径不是一个有效的文件: %v", filepath.Join(path, host, "cert.csr")))
	}

	issueStat, err := os.Stat(filepath.Join(path, host, "cert.issue"))
	if nil != err {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("检查证书失败: 无法找到证书文件: %v", filepath.Join(path, host, "cert.issue")))
		}

		return nil, errors.New(fmt.Sprintf("检查证书失败: %v", err))
	}

	if issueStat.IsDir() {
		return nil, errors.New(fmt.Sprintf("证书配置错误，此路径不是一个有效的文件: %v", filepath.Join(path, host, "cert.issue")))
	}

	cert, err := newCertificate(host)
	if nil != err {
		return nil, err
	}

	return cert, nil
}

// GetManager 获取证书管理器
func GetManager() *Manager {
	if nil == manager {
		logger.GetLogger().Errorln("调用了未初始化的证书管理器")

		panic("未初始化的证书管理器")
	}

	return manager
}

type Manager struct {
	ctx context.Context

	sync.RWMutex
	certs map[string]*Certificate
}

// RenewLoop 证书续期循环
// 该方法必须放在goroutine中运行
// 该方法在凌晨4-6点进行证书检查，并续期所有有效期在10天以内的证书
func (m *Manager) RenewLoop() {
	ticker := time.NewTicker(1 * time.Hour)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			hour := time.Now().Hour()
			if 4 > hour || 6 < hour {
				continue
			}

			m.RLock()
			for host, cert := range m.certs {
				if cert.GetExpireTime()-time.Now().Unix() > 10*24*3600 {
					continue
				}

				// 10天以内续期
				if err := m.renew(host); nil != err {
					logger.GetLogger().Errorln(err)
				}
			}
			m.RUnlock()
		}
	}
}

// renew 续期证书
// 该方法将会暂时关闭Caddy/V2ray服务
func (m *Manager) renew(host string) error {
	defer func() {
		if err := recover(); nil != err {
			logger.GetLogger().Errorf("Panic: %v", err)
		}
	}()

	// 检查Caddy
	caddyIsRunning, err := caddy.IsRunning()
	if nil != err {
		return errors.New(fmt.Sprintf("检查Caddy状态失败: %v", err))
	}

	if caddyIsRunning {
		if err := caddy.Stop(); nil != err {
			return errors.New(fmt.Sprintf("停止Caddy服务失败: %v", err))
		}
	}

	// 检查V2ray
	v2rayIsRunning, err := v2ray.IsRunning()
	if nil != err {
		return errors.New(fmt.Sprintf("检查V2ray状态失败: %v", err))
	}

	if v2rayIsRunning {
		if err := v2ray.Stop(); nil != err {
			return errors.New(fmt.Sprintf("停止V2ray服务失败: %v", err))
		}
	}

	httpIsAllow, err := util.CheckLocalPortIsAllow(80)
	if nil != err {
		return errors.New(fmt.Sprintf("检查端口状态失败: %v", err))
	}

	if !httpIsAllow {
		return errors.New("续期证书失败，80端口被占用")
	}

	httpsIsAllow, err := util.CheckLocalPortIsAllow(443)
	if nil != err {
		return errors.New(fmt.Sprintf("检查端口状态失败: %v", err))
	}

	if !httpsIsAllow {
		return errors.New("续期证书失败，443端口被占用")
	}

	if err := m.Renew(host, configurator.GetMainConfig().Email); nil != err {
		logger.GetLogger().Errorln(err)
	}

	// 启动Caddy服务
	if err := caddy.Start(); nil != err {
		return errors.New(fmt.Sprintf("启动Caddy服务失败: %v", err))
	}

	// 检查Caddy
	caddyIsRunning, err = caddy.IsRunning()
	if nil != err {
		return errors.New(fmt.Sprintf("检查Caddy状态失败: %v", err))
	}

	if !caddyIsRunning {
		return errors.New(fmt.Sprintf("启动Caddy服务失败: %v", err))
	}

	// 启动V2ray服务
	if err := v2ray.Start(); nil != err {
		return errors.New(fmt.Sprintf("启动V2ray服务失败: %v", err))
	}

	v2rayIsRunning, err = v2ray.IsRunning()
	if nil != err {
		return errors.New(fmt.Sprintf("检查V2ray状态失败: %v", err))
	}

	if !v2rayIsRunning {
		return errors.New(fmt.Sprintf("启动V2ray服务失败: %v", err))
	}

	return nil
}

// CheckExists 检查域名证书是否存在
func (m *Manager) CheckExists(host string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.certs[host]; ok {
		return true
	}

	return false
}

// GetCertificate 获取域名的证书
func (m *Manager) GetCertificate(host string) (*Certificate, error) {
	m.RLock()
	defer m.RUnlock()

	if cert, ok := m.certs[host]; ok {
		return cert, nil
	}

	return nil, errors.New("找不到该域名的证书")
}

// GetMustCertificate 获取域名的证书，该方法假定一定成功
func (m *Manager) GetMustCertificate(host string) *Certificate {
	m.RLock()
	defer m.RUnlock()

	return m.certs[host]
}

// IssueNew 申请新的证书
// 申请证书需要监听80与443端口，申请证书前需要关闭所有web服务，以免续期失败
func (m *Manager) IssueNew(host, email string) (*Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无法生成私钥: %v", err))
	}

	userEntry := &user{Email: email, key: privateKey}
	config := lego.NewConfig(userEntry)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("无法创建Acme客户端: %v", err))
	}

	if err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80")); err != nil {
		return nil, errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	if err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443")); err != nil {
		return nil, errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("注册用户失败: %v", err))
	}

	userEntry.Registration = reg

	request := certificate.ObtainRequest{Domains: []string{host}, Bundle: true}
	cert, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	if err := m.save(host, cert); nil != err {
		return nil, err
	}

	entry, err := newCertificate(host)
	if nil != err {
		return nil, err
	}

	return entry, nil
}

// Renew 续期证书
// 续期证书需要监听80与443端口，续期前需要关闭所有web服务器，以免续期失败
func (m *Manager) Renew(host, email string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return errors.New(fmt.Sprintf("无法生成私钥: %v", err))
	}

	lc, err := m.GetCertificate(host)
	if nil != err {
		return err
	}

	userEntry := &user{Email: email, key: privateKey}
	config := lego.NewConfig(userEntry)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return errors.New(fmt.Sprintf("无法创建Acme客户端: %v", err))
	}

	if err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80")); err != nil {
		return errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	if err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443")); err != nil {
		return errors.New(fmt.Sprintf("绑定端口监听失败: %v", err))
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return errors.New(fmt.Sprintf("注册用户失败: %v", err))
	}

	userEntry.Registration = reg

	request := certificate.Resource{Domain: host, PrivateKey: lc.GetPrivateKeyContent(), Certificate: lc.GetCertificateContent(), IssuerCertificate: lc.GetIssueCertificate(), CSR: lc.GetCsrContent()}
	cert, err := client.Certificate.Renew(request, true, false, "")
	if err != nil {
		return errors.New(fmt.Sprintf("续期证书失败: %v", err))
	}

	if err := m.save(host, cert); nil != err {
		return err
	}

	return nil
}

// save 保存证书
func (m *Manager) save(host string, cert *certificate.Resource) error {
	path := runtime.GetCertificatePath()

	if err := os.RemoveAll(filepath.Join(path, host)); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	if err := os.MkdirAll(filepath.Join(path, host), 0755); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存私钥
	keyFile, err := os.OpenFile(filepath.Join(path, host, "private.key"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer keyFile.Close()

	if _, err := keyFile.Write(cert.PrivateKey); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存证书
	certFile, err := os.OpenFile(filepath.Join(path, host, "cert.pem"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer certFile.Close()

	if _, err := certFile.Write(cert.Certificate); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存csr
	csrFile, err := os.OpenFile(filepath.Join(path, host, "cert.csr"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer csrFile.Close()

	if _, err := csrFile.Write(cert.CSR); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	// 保存IssuerCertificate
	issueFile, err := os.OpenFile(filepath.Join(path, host, "cert.issue"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}
	defer issueFile.Close()

	if _, err := issueFile.Write(cert.IssuerCertificate); nil != err {
		return errors.New(fmt.Sprintf("申请证书失败: %v", err))
	}

	return nil
}

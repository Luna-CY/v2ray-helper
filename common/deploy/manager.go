package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/software/nginx"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
	"github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/Luna-CY/v2ray-helper/dataservice"
	"sync"
	"time"
)

const (
	StagePreCheck = iota
	StageCertificateConfig
	StageCaddyConfig
	StageV2rayConfig
	StageFakeConfig
	StageConfigGenerate
)

const (
	StateRunning = iota + 1
	StateSuccess
	StateFailure
)

var manager *Manager

func GetManager() *Manager {
	if nil == manager {
		manager = new(Manager)
	}

	return manager
}

type Stage struct {
	Stage int // 阶段数，从0开始顺序表示：环境预检查、HTTPS证书、V2ray、伪装、Caddy、配置生成
	State int
}

type Manager struct {
	isRunning bool
	config    *Config

	sm     sync.RWMutex
	stages []Stage

	cloudreveAdminPassword string
}

// DeployServer 部署服务器
func (m *Manager) DeployServer(config *Config) error {
	if m.isRunning {
		return errors.New("已在运行一个部署任务，不能重复启动")
	}

	m.isRunning = true

	m.sm.Lock()
	defer m.sm.Unlock()
	m.stages = []Stage{}

	m.config = config
	m.cloudreveAdminPassword = ""

	go m.deploy()

	return nil
}

func (m *Manager) Clean() {
	m.sm.Lock()
	defer m.sm.Unlock()

	m.stages = []Stage{}
	m.cloudreveAdminPassword = ""
}

func (m *Manager) GetCloudreveAdminPassword() string {
	return m.cloudreveAdminPassword
}

func (m *Manager) GetStage() []Stage {
	if !m.isRunning {
		return nil
	}

	m.sm.RLock()
	defer m.sm.RUnlock()

	return m.stages
}

// deploy 部署服务器
func (m *Manager) deploy() {
	defer func() {
		if err := recover(); nil != err {
			logger.GetLogger().Errorf("Panic: %v", err)
		}
	}()

	defer func() {
		m.isRunning = false
	}()

	// 部署前检查
	m.switchStage(StagePreCheck, StateRunning)
	if err := m.preCheck(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StagePreCheck, StateFailure)

		return
	}
	m.switchStage(StagePreCheck, StateSuccess)

	// 部署证书
	m.switchStage(StageCertificateConfig, StateRunning)
	if err := m.deployCertificate(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StageCertificateConfig, StateFailure)

		return
	}
	m.switchStage(StageCertificateConfig, StateSuccess)

	// 部署V2ray
	m.switchStage(StageV2rayConfig, StateRunning)
	if err := m.deployV2ray(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StageV2rayConfig, StateFailure)

		return
	}
	m.switchStage(StageV2rayConfig, StateSuccess)

	// 部署伪装站点
	m.switchStage(StageFakeConfig, StateRunning)
	if err := m.deployFakeWebServer(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StageFakeConfig, StateFailure)

		return
	}
	m.switchStage(StageFakeConfig, StateSuccess)

	// 部署Caddy
	m.switchStage(StageCaddyConfig, StateRunning)
	if err := m.deployCaddy(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StageCaddyConfig, StateFailure)

		return
	}
	m.switchStage(StageCaddyConfig, StateSuccess)

	// 重新生成客户端配置
	m.switchStage(StageConfigGenerate, StateRunning)
	if err := m.generateClientConfig(); nil != err {
		logger.GetLogger().Errorln(err)
		m.switchStage(StageConfigGenerate, StateFailure)
	}
	m.switchStage(StageConfigGenerate, StateSuccess)
}

// preCheck 预检查
func (m *Manager) preCheck() error {
	// 如果有Nginx服务器并且已启动，那么停止Nginx，否则无法申请证书且Caddy无法启动
	nginxIsRunning, err := nginx.IsRunning()
	if nil != err {
		return err
	}

	if nginxIsRunning {
		if err := nginx.Stop(); nil != err {
			return err
		}
	}

	if err := nginx.Disable(); nil != err {
		return err
	}

	return nil
}

// deployCertificate 部署证书
func (m *Manager) deployCertificate() error {
	if "" == m.config.HttpsHost {
		return nil
	}

	if !certificate.GetManager().CheckExists(m.config.HttpsHost) {
		_, err := certificate.GetManager().IssueNew(m.config.HttpsHost, configurator.GetMainConfig().Email)
		if nil != err {
			return err
		}

		return nil
	}

	return nil
}

// deployV2ray 部署V2ray服务
func (m *Manager) deployV2ray() error {
	if err := v2ray.Stop(); nil != err {
		return err
	}

	if err := v2ray.InstallLastRelease(); nil != err {
		return err
	}

	if err := v2ray.SetConfig(v2ray.ConfigPath, m.config.V2rayConfig); nil != err {
		return err
	}

	// 启动V2ray服务
	if err := v2ray.Start(); nil != err {
		return err
	}

	// 设为开机自启
	if err := v2ray.Enable(); nil != err {
		return err
	}

	return nil
}

// generateClientConfig 生成客户端配置
func (m *Manager) generateClientConfig() error {
	tcp, err := json.Marshal(m.config.V2rayConfig.Tcp)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	webSocket, err := json.Marshal(m.config.V2rayConfig.WebSocket)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	kcp, err := json.Marshal(m.config.V2rayConfig.Kcp)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	http2, err := json.Marshal(m.config.V2rayConfig.Http2)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	tcpString := string(tcp)
	webSocketString := string(webSocket)
	kcpString := string(kcp)
	http2String := string(http2)

	host := m.config.HttpsHost
	if "" == m.config.HttpsHost {
		ip, err := util.GetPublicIpv4()
		if nil != err {
			return errors.New(fmt.Sprintf("获取本机IP失败: %v", err))
		}

		host = ip
	}

	port := 80
	if "" != m.config.HttpsHost {
		port = 443
	}

	if v2ray.TransportTypeTcp == m.config.V2rayConfig.TransportType || v2ray.TransportTypeKcp == m.config.V2rayConfig.TransportType {
		port = m.config.V2rayConfig.V2rayPort
	}

	useTls := 0
	if "" != m.config.HttpsHost {
		useTls = 1
	}

	one := 1

	for _, client := range m.config.V2rayConfig.Clients {
		endpoint := model.V2rayEndpoint{
			Cloud:         &one,
			Endpoint:      &one,
			Host:          &host,
			Port:          &port,
			UserId:        &client.UserId,
			AlterId:       &client.AlterId,
			UseTls:        &useTls,
			TransportType: &m.config.V2rayConfig.TransportType,
			Tcp:           &tcpString,
			WebSocket:     &webSocketString,
			Kcp:           &kcpString,
			Http2:         &http2String,
		}

		ct := time.Now().Unix()
		endpoint.CreateTime = &ct
		endpoint.Deleted = util.NewFalsePtr()

		if err := dataservice.GetBaseService().Create(&endpoint); nil != err {
			return errors.New(fmt.Sprintf("保存数据失败: %v", err))
		}
	}

	return nil
}

func (m *Manager) switchStage(stage int, state int) {
	m.sm.Lock()
	defer m.sm.Unlock()

	m.stages = append(m.stages, Stage{Stage: stage, State: state})
}

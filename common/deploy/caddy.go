package deploy

import (
	"github.com/Luna-CY/v2ray-helper/common/software/caddy"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
)

// deployCaddy 部署Caddy服务
func (m *Manager) deployCaddy() error {
	if err := caddy.Stop(); nil != err {
		return err
	}

	if err := caddy.InstallLastRelease(); nil != err {
		return err
	}

	port := caddy.PortHttp
	if "" != m.config.HttpsHost {
		port = caddy.PortHttps
	}

	enableCloudreve := false
	if nil != m.config.FakeConfig && FakeTypeCloudreve == m.config.FakeConfig.FakeType {
		enableCloudreve = true
	}

	enableHttp2 := false
	if v2ray.TransportTypeHttp2 == m.config.V2rayConfig.TransportType {
		enableHttp2 = true
	}

	if err := caddy.SetConfigToSystem(m.config.HttpsHost, port, m.config.V2rayConfig.V2rayPort, m.config.V2rayConfig.WebSocket.Path, "" != m.config.HttpsHost, enableCloudreve, enableHttp2); nil != err {
		return err
	}

	if err := caddy.Start(); nil != err {
		return nil
	}

	if err := caddy.Enable(); nil != err {
		return nil
	}

	return nil
}

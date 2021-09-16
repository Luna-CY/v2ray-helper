package deploy

import (
	"github.com/Luna-CY/v2ray-helper/common/software/aria2"
	"github.com/Luna-CY/v2ray-helper/common/software/cloudreve"
)

func (m *Manager) deployFakeWebServer() error {
	if nil == m.config.FakeConfig || FakeTypeNone == m.config.FakeConfig.FakeType {
		return nil
	}

	if FakeTypeCloudreve == m.config.FakeConfig.FakeType && nil != m.config.FakeConfig.CloudreveConfig {
		if err := cloudreve.Stop(); nil != err {
			return err
		}

		if err := cloudreve.InstallLastRelease(); nil != err {
			return err
		}

		if err := cloudreve.Start(); nil != err {
			return err
		}

		if err := cloudreve.Enable(); nil != err {
			return err
		}

		if m.config.FakeConfig.CloudreveConfig.EnableAria2 {
			if err := m.deployAria2(); nil != err {
				return err
			}
		}

		if m.config.FakeConfig.CloudreveConfig.ResetAdminPassword {
			password, err := cloudreve.ResetAdminPassword()
			if nil != err {
				return err
			}

			m.cloudreveAdminPassword = password
		}
	}

	return nil
}

// deployAria2 部署Aria2
func (m *Manager) deployAria2() error {
	if err := aria2.InstallToSystem(); nil != err {
		return err
	}

	if err := aria2.Start(); nil != err {
		return err
	}

	if err := aria2.Enable(); nil != err {
		return err
	}

	if err := cloudreve.SetAria2(cloudreve.DefaultDbPath, "http://127.0.0.1:6800", aria2.DefaultToken, cloudreve.Aria2TempPath); nil != err {
		return err
	}

	return nil
}

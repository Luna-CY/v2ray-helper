package deploy

import (
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
)

type Config struct {
	HttpsHost   string
	V2rayConfig *v2ray.Config
	FakeConfig  *FakeWebServerConfig
}

const (
	FakeTypeNone = iota
	FakeTypeCloudreve
)

type FakeWebServerConfig struct {
	FakeType        int
	CloudreveConfig *CloudreveConfig
}

type CloudreveConfig struct {
	EnableAria2        bool
	ResetAdminPassword bool
}

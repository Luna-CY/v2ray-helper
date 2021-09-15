package caddy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"os"
	"path/filepath"
)

const (
	PortHttp  = 80
	PortHttps = 443
)

// SetConfig 添加Caddy的配置
func SetConfig(configPath, host string, port, v2rayPort int, path string, cloudreve, http2 bool) error {
	builder := bytes.Buffer{}
	builder.WriteString(fmt.Sprintf("%v:%v {\n", host, port))

	if PortHttps == port {
		certFilePath := filepath.Join(runtime.GetCertificatePath(), host, "cert.pem")
		keyFilePath := filepath.Join(runtime.GetCertificatePath(), host, "private.key")

		builder.WriteString(fmt.Sprintf("    tls %v %v\n", certFilePath, keyFilePath))
	}

	if http2 {
		builder.WriteString(fmt.Sprintf("    reverse_proxy %v 127.0.0.1:%v {\n", path, v2rayPort))
		builder.WriteString("        transport http {\n")
		builder.WriteString("            versions h2c\n")
		builder.WriteString("        }\n")
		builder.WriteString("    }\n")
	} else {
		builder.WriteString(fmt.Sprintf("    reverse_proxy %v 127.0.0.1:%v\n", path, v2rayPort))
	}

	if cloudreve {
		builder.WriteString("    reverse_proxy 127.0.0.1:5212\n")
	}
	builder.WriteString("}")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("配置Caddy失败: %v", err))
	}

	cf, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
	}
	defer cf.Close()

	if _, err := cf.Write(builder.Bytes()); nil != err {
		return errors.New(fmt.Sprintf("无法写入配置文件: %v", err))
	}

	return nil
}

// SetConfigToSystem 设置Caddy
func SetConfigToSystem(host string, port, v2rayPort int, path string, cloudreve, http2 bool) error {
	return SetConfig(ConfigPath, host, port, v2rayPort, path, cloudreve, http2)
}

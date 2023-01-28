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
func SetConfig(configPath, host string, listenPort, targetPort int, path string, https, cloudreve, http2 bool) error {
	builder := bytes.Buffer{}
	builder.WriteString(fmt.Sprintf("%v:%v {\n", host, listenPort))

	if PortHttps == listenPort || https || http2 {
		certFilePath := filepath.Join(runtime.GetAcmeCertificatePath(), host, "cert.pem")
		keyFilePath := filepath.Join(runtime.GetAcmeCertificatePath(), host, "private.key")

		builder.WriteString(fmt.Sprintf("    tls %v %v\n", certFilePath, keyFilePath))
	}

	if http2 {
		builder.WriteString(fmt.Sprintf("    reverse_proxy %v 127.0.0.1:%v {\n", path, targetPort))
		builder.WriteString("        transport http {\n")
		builder.WriteString("            versions h2c\n")
		builder.WriteString("        }\n")
		builder.WriteString("    }\n")
	} else {
		builder.WriteString(fmt.Sprintf("    reverse_proxy %v 127.0.0.1:%v\n", path, targetPort))
	}

	if cloudreve {
		builder.WriteString("    reverse_proxy 127.0.0.1:5212\n")
	}
	builder.WriteString("}")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("配置Caddy失败: %v", err))
	}
	if err := os.RemoveAll(configPath); nil != err {
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
func SetConfigToSystem(host string, listenPort, targetPort int, path string, https, cloudreve, http2 bool) error {
	configPath := filepath.Join(filepath.Dir(ConfigPath), "sites-enabled", fmt.Sprintf("%v.%v", host, listenPort))

	return SetConfig(configPath, host, listenPort, targetPort, path, https, cloudreve, http2)
}

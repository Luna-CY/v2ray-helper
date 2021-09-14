package v2ray

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"os"
	"path/filepath"
	"strings"
)

const ConfigPath = "/usr/local/etc/v2ray/config.json"

const (
	TransportTypeTcp = iota + 1
	TransportTypeWebSocket
	TransportTypeKcp
	TransportTypeHttp2
)

const (
	TcpTypeNone = "none"
	TcpTypeHttp = "http"
)

// Config 配置结构
type Config struct {
	Clients []ConfigClient `json:"clients"`

	V2rayPort     int `json:"v2ray_port"`
	TransportType int `json:"transport_type"`

	Tcp struct {
		Type    string `json:"type"`
		Request struct {
			Version string         `json:"version"`
			Method  string         `json:"method"`
			Path    string         `json:"path"`
			Headers []ConfigHeader `json:"headers"`
		} `json:"request"`
		Response struct {
			Version string         `json:"version"`
			Status  string         `json:"status"`
			Reason  string         `json:"reason"`
			Headers []ConfigHeader `json:"headers"`
		} `json:"response"`
	} `json:"tcp"`
	WebSocket struct {
		Path    string         `json:"path"`
		Headers []ConfigHeader `json:"headers"`
	} `json:"web_socket"`
	Kcp struct {
		Type             string `json:"type"`
		Mtu              int    `json:"mtu"`
		Tti              int    `json:"tti"`
		UpLinkCapacity   int    `json:"uplink_capacity"`
		DownLinkCapacity int    `json:"downlink_capacity"`
		Congestion       bool   `json:"congestion"`
		ReadBufferSize   int    `json:"read_buffer_size"`
		WriteBufferSize  int    `json:"write_buffer_size"`
	} `json:"kcp"`
	Http2 struct {
		Host string `json:"host"`
		Path string `json:"path"`
	} `json:"http2"`

	UseTls  bool
	TlsHost string
	TlsKey  []byte
	TlsCert []byte
}

type ConfigClient struct {
	UserId  string `json:"user_id"`
	AlterId int    `json:"alter_id"`
}

type ConfigHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// vConfig V2ray配置结构
type vConfig struct {
	Inbounds  []vConfigInbound  `json:"inbounds"`
	Outbounds []vConfigOutbound `json:"outbounds"`
}

type vConfigInbound struct {
	Listen   string `json:"listen"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Settings struct {
		Clients []vConfigInboundClient `json:"clients"`
	} `json:"settings"`
	StreamSettings struct {
		Network      string                         `json:"network"`
		Security     string                         `json:"security"`
		TlsSettings  *vConfigInboundTls             `json:"tlsSettings,omitempty"`
		TcpSettings  *vConfigInboundStreamTcp       `json:"tcpSettings,omitempty"`
		WsSettings   *vConfigInboundStreamWebSocket `json:"wsSettings,omitempty"`
		KcpSettings  *vConfigInboundStreamKcp       `json:"kcpSettings,omitempty"`
		HttpSettings *vConfigInboundStreamHttp      `json:"httpSettings,omitempty"`
	} `json:"streamSettings"`
}

type vConfigInboundClient struct {
	Id      string `json:"id"`
	AlterId int    `json:"alterId"`
}

type vConfigInboundTls struct {
	ServerName        string                         `json:"serverName"`
	AllowInsecure     bool                           `json:"allowInsecure"`
	Alpn              []string                       `json:"alpn"`
	Certificates      []vConfigInboundTlsCertificate `json:"certificates"`
	DisableSystemRoot bool                           `json:"disableSystemRoot"`
}

type vConfigInboundTlsCertificate struct {
	Usage           string   `json:"usage"`
	CertificateFile string   `json:"certificateFile"`
	KeyFile         string   `json:"keyFile"`
	Certificate     []string `json:"certificate"`
	Key             []string `json:"key"`
}

type vConfigInboundStreamTcp struct {
	Header struct {
		Type     string                           `json:"type"`
		Request  *vConfigInboundStreamTcpRequest  `json:"request,omitempty"`
		Response *vConfigInboundStreamTcpResponse `json:"response,omitempty"`
	} `json:"header"`
}

type vConfigInboundStreamTcpRequest struct {
	Version string              `json:"version"`
	Method  string              `json:"method"`
	Path    []string            `json:"path"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type vConfigInboundStreamTcpResponse struct {
	Version string              `json:"version"`
	Status  string              `json:"status"`
	Reason  string              `json:"reason"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type vConfigInboundStreamWebSocket struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

type vConfigInboundStreamKcp struct {
	Header struct {
		Type string `json:"type"`
	} `json:"header"`
	Mtu              int  `json:"mtu"`
	Tti              int  `json:"tti"`
	UpLinkCapacity   int  `json:"uplinkCapacity"`
	DownLinkCapacity int  `json:"downlinkCapacity"`
	Congestion       bool `json:"congestion"`
	ReadBufferSize   int  `json:"readBufferSize"`
	WriteBufferSize  int  `json:"writeBufferSize"`
}

type vConfigInboundStreamHttp struct {
	Host []string `json:"host"`
	Path string   `json:"path"`
}

type vConfigOutbound struct {
	Protocol string   `json:"protocol"`
	Settings struct{} `json:"settings"`
}

// SetConfig 设置V2ray配置
func SetConfig(configPath string, config *Config) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); nil != err {
		return errors.New(fmt.Sprintf("配置V2ray失败: %v", err))
	}

	inbound := vConfigInbound{
		Listen:   "127.0.0.1",
		Port:     config.V2rayPort,
		Protocol: "vmess",
	}

	for _, client := range config.Clients {
		inbound.Settings.Clients = append(inbound.Settings.Clients, vConfigInboundClient{
			Id:      client.UserId,
			AlterId: client.AlterId,
		})
	}

	inbound.StreamSettings.Security = "none"

	cert := vConfigInboundTlsCertificate{
		Usage: "encipherment",
	}

	if 0 != len(config.TlsKey) && 0 != len(config.TlsCert) {
		cert.Key = strings.Split(string(config.TlsKey), "\n")
		cert.Certificate = strings.Split(string(config.TlsCert), "\n")
	} else {
		cert.KeyFile = filepath.Join(runtime.GetRootPath(), certificate.CertDirName, config.TlsHost, "private.key")
		cert.CertificateFile = filepath.Join(runtime.GetRootPath(), certificate.CertDirName, config.TlsHost, "cert.pem")
	}

	switch config.TransportType {
	case TransportTypeTcp:
		// TCP协议需要V2ray来监听端口
		inbound.Listen = "0.0.0.0"

		inbound.StreamSettings.Network = "tcp"
		inbound.StreamSettings.TcpSettings = &vConfigInboundStreamTcp{}
		inbound.StreamSettings.TcpSettings.Header.Type = config.Tcp.Type
		if TcpTypeHttp == config.Tcp.Type {
			inbound.StreamSettings.TcpSettings.Header.Request = &vConfigInboundStreamTcpRequest{}
			inbound.StreamSettings.TcpSettings.Header.Request.Version = config.Tcp.Request.Version
			inbound.StreamSettings.TcpSettings.Header.Request.Method = config.Tcp.Request.Method
			inbound.StreamSettings.TcpSettings.Header.Request.Path = strings.Split(config.Tcp.Request.Path, ",")

			if 0 < len(config.Tcp.Request.Headers) {
				inbound.StreamSettings.TcpSettings.Header.Request.Headers = make(map[string][]string)
				for _, header := range config.Tcp.Request.Headers {
					inbound.StreamSettings.TcpSettings.Header.Request.Headers[header.Key] = strings.Split(header.Value, ";;;")
				}
			}

			inbound.StreamSettings.TcpSettings.Header.Response = &vConfigInboundStreamTcpResponse{}
			inbound.StreamSettings.TcpSettings.Header.Response.Version = config.Tcp.Response.Version
			inbound.StreamSettings.TcpSettings.Header.Response.Status = config.Tcp.Response.Status
			inbound.StreamSettings.TcpSettings.Header.Response.Reason = config.Tcp.Response.Reason

			if 0 < len(config.Tcp.Response.Headers) {
				inbound.StreamSettings.TcpSettings.Header.Response.Headers = make(map[string][]string)
				for _, header := range config.Tcp.Response.Headers {
					inbound.StreamSettings.TcpSettings.Header.Response.Headers[header.Key] = strings.Split(header.Value, ";;;")
				}
			}
		}

		// 启用TLS
		if config.UseTls {
			inbound.StreamSettings.Security = "tls"
			inbound.StreamSettings.TlsSettings.ServerName = config.TlsHost
			inbound.StreamSettings.TlsSettings.Certificates = append(inbound.StreamSettings.TlsSettings.Certificates, cert)
		}
	case TransportTypeWebSocket:
		inbound.StreamSettings.Network = "ws"
		inbound.StreamSettings.WsSettings = &vConfigInboundStreamWebSocket{}
		inbound.StreamSettings.WsSettings.Path = config.WebSocket.Path
		inbound.StreamSettings.WsSettings.Headers = make(map[string]string)
		if 0 < len(config.WebSocket.Headers) {
			for _, header := range config.WebSocket.Headers {
				inbound.StreamSettings.WsSettings.Headers[header.Key] = header.Value
			}
		}
	case TransportTypeKcp:
		// KCP协议需要V2ray来监听端口
		inbound.Listen = "0.0.0.0"

		inbound.StreamSettings.Network = "kcp"
		inbound.StreamSettings.KcpSettings = &vConfigInboundStreamKcp{}
		inbound.StreamSettings.KcpSettings.Mtu = config.Kcp.Mtu
		inbound.StreamSettings.KcpSettings.Tti = config.Kcp.Tti
		inbound.StreamSettings.KcpSettings.Congestion = config.Kcp.Congestion
		inbound.StreamSettings.KcpSettings.UpLinkCapacity = config.Kcp.UpLinkCapacity
		inbound.StreamSettings.KcpSettings.DownLinkCapacity = config.Kcp.DownLinkCapacity
		inbound.StreamSettings.KcpSettings.ReadBufferSize = config.Kcp.ReadBufferSize
		inbound.StreamSettings.KcpSettings.WriteBufferSize = config.Kcp.WriteBufferSize
		inbound.StreamSettings.KcpSettings.Header.Type = config.Kcp.Type
	case TransportTypeHttp2:
		inbound.StreamSettings.Network = "http"
		inbound.StreamSettings.HttpSettings = &vConfigInboundStreamHttp{}
		inbound.StreamSettings.HttpSettings.Host = strings.Split(config.Http2.Host, ",")
		inbound.StreamSettings.HttpSettings.Path = config.Http2.Path

		// 启用TLS
		if config.UseTls {
			inbound.StreamSettings.Security = "tls"
			inbound.StreamSettings.TlsSettings.ServerName = config.TlsHost
			inbound.StreamSettings.TlsSettings.Certificates = append(inbound.StreamSettings.TlsSettings.Certificates, cert)
		}
	default:
		return errors.New(fmt.Sprintf("未受支持的传输方式: %v", config.TransportType))
	}

	outbound := vConfigOutbound{
		Protocol: "freedom",
		Settings: struct{}{},
	}

	vc := vConfig{
		Inbounds:  []vConfigInbound{inbound},
		Outbounds: []vConfigOutbound{outbound},
	}

	content, err := json.Marshal(vc)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	cf, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
	}
	defer cf.Close()

	if _, err := cf.Write(content); nil != err {
		return errors.New(fmt.Sprintf("写入配置失败: %v", err))
	}

	return nil
}

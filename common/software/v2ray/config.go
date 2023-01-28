package v2ray

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// VConfig V2ray配置结构
type VConfig struct {
	Inbounds  []VConfigInbound  `json:"inbounds"`
	Outbounds []VConfigOutbound `json:"outbounds"`
}

type VConfigInbound struct {
	Listen   string `json:"listen"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Settings struct {
		Clients []VConfigInboundClient `json:"clients"`
	} `json:"settings"`
	StreamSettings struct {
		Network      string                         `json:"network"`
		Security     string                         `json:"security"`
		TlsSettings  *VConfigInboundTls             `json:"tlsSettings,omitempty"`
		TcpSettings  *VConfigInboundStreamTcp       `json:"tcpSettings,omitempty"`
		WsSettings   *VConfigInboundStreamWebSocket `json:"wsSettings,omitempty"`
		KcpSettings  *VConfigInboundStreamKcp       `json:"kcpSettings,omitempty"`
		HttpSettings *VConfigInboundStreamHttp      `json:"httpSettings,omitempty"`
	} `json:"streamSettings"`
}

type VConfigInboundClient struct {
	Id      string `json:"id"`
	AlterId int    `json:"alterId"`
}

type VConfigInboundTls struct {
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

type VConfigInboundStreamTcp struct {
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

type VConfigInboundStreamWebSocket struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

type VConfigInboundStreamKcp struct {
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

type VConfigInboundStreamHttp struct {
	Host []string `json:"host"`
	Path string   `json:"path"`
}

type VConfigOutbound struct {
	Protocol string   `json:"protocol"`
	Settings struct{} `json:"settings"`
}

// GetConfig 读取配置文件
func GetConfig(configPath string) (*VConfig, error) {
	file, err := os.Open(configPath)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法打开文件: %v", err))
	}

	content, err := io.ReadAll(file)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("无法读取配置文件: %v", err))
	}

	var vc VConfig
	if err := json.Unmarshal(content, &vc); nil != err {
		return nil, errors.New(fmt.Sprintf("无法解析配置文件: %v", err))
	}

	return &vc, nil
}

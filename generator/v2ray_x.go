package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
)

// GenerateV2rayXContent 生成V2rayX的配置内容
func GenerateV2rayXContent(endpoint model.V2rayEndpoint) (string, error) {
	v2rayXConfig := V2rayXConfig{}

	outbound := V2rayXOutbound{
		SendThrough: "0.0.0.0",
		Protocol:    "vmess",
		Tag:         fmt.Sprintf("%v:%v", *endpoint.Host, *endpoint.Port),
	}

	outbound.Mux.Enabled = false
	outbound.Mux.Concurrency = 8

	vnext := V2rayXOutboundVNext{
		Address: *endpoint.Host,
		Port:    *endpoint.Port,
	}

	user := V2rayXOutboundVNextUser{
		Id:       *endpoint.UserId,
		AlterId:  *endpoint.AlterId,
		Security: "auto",
		Level:    *endpoint.Level,
	}

	vnext.Users = append(vnext.Users, user)
	outbound.Settings.Vnext = append(outbound.Settings.Vnext, vnext)

	outbound.StreamSettings.Network = "tcp"
	outbound.StreamSettings.Security = "none"

	outbound.StreamSettings.TcpSettings.Header.Type = "none"

	outbound.StreamSettings.QuicSettings.Security = "none"
	outbound.StreamSettings.QuicSettings.Header.Type = "none"

	outbound.StreamSettings.TlsSettings.AllowInsecure = false
	outbound.StreamSettings.TlsSettings.AllowInsecureCiphers = false
	outbound.StreamSettings.TlsSettings.Alpn = []string{"http/1.1"}

	outbound.StreamSettings.KcpSettings.Header.Type = "none"
	outbound.StreamSettings.KcpSettings.Mtu = 1350
	outbound.StreamSettings.KcpSettings.Tti = 20
	outbound.StreamSettings.KcpSettings.UplinkCapacity = 5
	outbound.StreamSettings.KcpSettings.WriteBufferSize = 1
	outbound.StreamSettings.KcpSettings.ReadBufferSize = 1
	outbound.StreamSettings.KcpSettings.DownlinkCapacity = 20

	if 1 == *endpoint.TransportType {
		outbound.StreamSettings.Network = "ws"
		outbound.StreamSettings.Security = "tls"

		if err := json.Unmarshal([]byte(*endpoint.WebSocket), &outbound.StreamSettings.WebSocket); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}
	}

	v2rayXConfig.Outbounds = append(v2rayXConfig.Outbounds, outbound)
	v2rayXConfig.Routes = make([]V2rayXRoute, 0)

	content, err := json.Marshal(v2rayXConfig)
	if nil != err {
		return "", errors.New("序列化数据失败，稍后重试一下吧，或者联系管理员")
	}

	return string(content), nil
}

type V2rayXConfig struct {
	Outbounds []V2rayXOutbound `json:"outbounds"`
	Routes    []V2rayXRoute    `json:"routings"`
}

type V2rayXOutbound struct {
	SendThrough string `json:"sendThrough"`
	Mux         struct {
		Enabled     bool `json:"enabled"`
		Concurrency int  `json:"concurrency"`
	} `json:"mux"`
	Protocol string `json:"protocol"`
	Settings struct {
		Vnext []V2rayXOutboundVNext `json:"vnext"`
	} `json:"settings"`
	Tag            string `json:"tag"`
	StreamSettings struct {
		Network   string `json:"network"`
		Security  string `json:"security"`
		WebSocket struct {
			Path    string   `json:"path"`
			Headers struct{} `json:"headers"`
		} `json:"wsSettings"`
		QuicSettings struct {
			Key      string `json:"key"`
			Security string `json:"security"`
			Header   struct {
				Type string `json:"type"`
			} `json:"header"`
		} `json:"quicSettings"`
		TlsSettings struct {
			AllowInsecure        bool     `json:"allowInsecure"`
			Alpn                 []string `json:"alpn"`
			ServerName           string   `json:"serverName"`
			AllowInsecureCiphers bool     `json:"allowInsecureCiphers"`
		} `json:"tlsSettings"`
		HttpSettings struct {
			Path string `json:"path"`
		} `json:"httpSettings"`
		KcpSettings struct {
			Header struct {
				Type string `json:"type"`
			} `json:"header"`
			Mtu              int  `json:"mtu"`
			Congestion       bool `json:"congestion"`
			Tti              int  `json:"tti"`
			UplinkCapacity   int  `json:"uplinkCapacity"`
			WriteBufferSize  int  `json:"writeBufferSize"`
			ReadBufferSize   int  `json:"readBufferSize"`
			DownlinkCapacity int  `json:"downlinkCapacity"`
		} `json:"kcpSettings"`
		TcpSettings struct {
			Header struct {
				Type string `json:"type"`
			} `json:"header"`
		} `json:"tcpSettings"`
	} `json:"streamSettings"`
}

type V2rayXOutboundVNext struct {
	Address string                    `json:"address"`
	Port    int                       `json:"port"`
	Users   []V2rayXOutboundVNextUser `json:"users"`
}

type V2rayXOutboundVNextUser struct {
	Id       string `json:"id"`
	AlterId  int    `json:"alter_id"`
	Security string `json:"security"`
	Level    int    `json:"level"`
}

type V2rayXRoute struct {
}

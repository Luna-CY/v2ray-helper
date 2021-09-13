package generator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"strconv"
)

// GenerateVMessShareLink 生成通用分享连接
func GenerateVMessShareLink(config model.V2rayEndpoint) (string, error) {
	vMess := VMessShareLinkProtocol{}
	vMess.V = "2"
	vMess.Scy = "auto"
	vMess.Sni = *config.Sni

	vMess.Add = *config.Host
	vMess.Port = strconv.Itoa(*config.Port)
	vMess.Id = *config.UserId
	vMess.Aid = strconv.Itoa(*config.AlterId)
	vMess.Ps = fmt.Sprintf("%v:%v", *config.Host, *config.Port)

	vMess.Tls = ""
	if 1 == *config.UseTls {
		vMess.Tls = "tls"
		vMess.Host = *config.Host
	}

	if model.V2rayEndpointTransportTypeTcp == *config.TransportType {
		vMess.Net = "tcp"

		var tcp model.V2rayEndpointTcp
		if err := json.Unmarshal([]byte(*config.Tcp), &tcp); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Type = tcp.Type
	}

	if model.V2rayEndpointTransportTypeWebSocket == *config.TransportType {
		vMess.Net = "ws"
		vMess.Type = ""

		var webSocket model.V2rayEndpointWebSocket
		if err := json.Unmarshal([]byte(*config.WebSocket), &webSocket); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Path = webSocket.Path
	}

	if model.V2rayEndpointTransportTypeKcp == *config.TransportType {
		vMess.Net = "kcp"

		var kcp model.V2rayEndpointKcp
		if err := json.Unmarshal([]byte(*config.Kcp), &kcp); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Type = kcp.Type
	}

	if model.V2rayEndpointTransportTypeHttp2 == *config.TransportType {
		vMess.Net = "h2"
		vMess.Type = ""

		var http2 model.V2rayEndpointHttp2
		if err := json.Unmarshal([]byte(*config.Http2), &http2); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Host = http2.Host
		vMess.Path = http2.Path
	}

	content, err := json.Marshal(vMess)
	if nil != err {
		return "", errors.New("序列化数据失败，稍后重试一下吧，或者联系管理员")
	}

	return fmt.Sprintf("vmess://%v", base64.StdEncoding.EncodeToString(content)), nil
}

type VMessShareLinkProtocol struct {
	Add  string `json:"add"`
	Port string `json:"port"`
	Aid  string `json:"aid"`
	Host string `json:"host"`
	Id   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Ps   string `json:"ps"`
	Scy  string `json:"scy"`
	Sni  string `json:"sni"`
	Tls  string `json:"tls"`
	Type string `json:"type"`
	V    string `json:"v"`
}

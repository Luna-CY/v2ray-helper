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
func GenerateVMessShareLink(endpoint model.V2rayEndpoint) (string, error) {
	vMess := VMessShareLinkProtocol{}
	vMess.V = "2"
	vMess.Scy = "auto"
	vMess.Sni = *endpoint.Sni

	vMess.Add = *endpoint.Host
	vMess.Port = strconv.Itoa(*endpoint.Port)
	vMess.Id = *endpoint.UserId
	vMess.Aid = strconv.Itoa(*endpoint.AlterId)
	vMess.Ps = fmt.Sprintf("%v:%v", *endpoint.Host, *endpoint.Port)

	vMess.Tls = ""
	if 1 == *endpoint.UseTls {
		vMess.Tls = "tls"
	}

	if model.V2rayEndpointTransportTypeTcp == *endpoint.TransportType {
		vMess.Net = "tcp"

		var tcp map[string]interface{}
		if err := json.Unmarshal([]byte(*endpoint.Tcp), &tcp); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := tcp["type"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Type = tcp["type"].(string)
	}

	if model.V2rayEndpointTransportTypeWebSocket == *endpoint.TransportType {
		vMess.Net = "ws"
		vMess.Type = ""

		var webSocket map[string]interface{}
		if err := json.Unmarshal([]byte(*endpoint.WebSocket), &webSocket); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := webSocket["path"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Path = webSocket["path"].(string)
	}

	if model.V2rayEndpointTransportTypeKcp == *endpoint.TransportType {
		vMess.Net = "kcp"

		var kcp map[string]interface{}
		if err := json.Unmarshal([]byte(*endpoint.Kcp), &kcp); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := kcp["type"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Type = kcp["type"].(string)
	}

	if model.V2rayEndpointTransportTypeHttp2 == *endpoint.TransportType {
		vMess.Net = "h2"
		vMess.Type = ""

		var http2 map[string]interface{}
		if err := json.Unmarshal([]byte(*endpoint.Http2), &http2); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := http2["host"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := http2["path"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		vMess.Host = http2["host"].(string)
		vMess.Path = http2["path"].(string)
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

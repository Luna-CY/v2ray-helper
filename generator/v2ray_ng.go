package generator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"strconv"
)

func GenerateV2rayNGContent(endpoint model.V2rayEndpoint) (string, error) {
	v2rayNGConfig := V2rayNGConfig{}

	v2rayNGConfig.Add = *endpoint.Host
	v2rayNGConfig.Port = strconv.Itoa(*endpoint.Port)
	v2rayNGConfig.Aid = strconv.Itoa(*endpoint.AlterId)

	v2rayNGConfig.Id = *endpoint.UserId
	v2rayNGConfig.Ps = fmt.Sprintf("%v:%v", *endpoint.Host, *endpoint.Port)

	v2rayNGConfig.V = "2"
	v2rayNGConfig.Scy = "auto"
	v2rayNGConfig.Net = "tcp"
	v2rayNGConfig.Type = "none"

	if model.V2rayEndpointTransportTypeWebSocket == *endpoint.TransportType {
		v2rayNGConfig.Net = "ws"
		v2rayNGConfig.Tls = "tls"
		v2rayNGConfig.Type = ""

		var webSocket map[string]interface{}
		if err := json.Unmarshal([]byte(*endpoint.WebSocket), &webSocket); nil != err {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		if _, ok := webSocket["path"]; !ok {
			return "", errors.New("配置错误，稍后重试一下吧，或者联系管理员")
		}

		v2rayNGConfig.Path = webSocket["path"].(string)
	}

	content, err := json.Marshal(v2rayNGConfig)
	if nil != err {
		return "", errors.New("序列化数据失败，稍后重试一下吧，或者联系管理员")
	}

	return fmt.Sprintf("vmess://%v", base64.StdEncoding.EncodeToString(content)), nil
}

type V2rayNGConfig struct {
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

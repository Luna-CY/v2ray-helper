package generator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// GenerateVMessShareLink 生成通用分享连接
func GenerateVMessShareLink() (string, error) {
	vMess := VMessShareLinkProtocol{}
	vMess.V = "2"
	vMess.Scy = "auto"

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

package controller

import "github.com/gin-gonic/gin"

type V2rayServerDeployForm struct {
	ServerType int `json:"server_type"`
	ServerIp string `json:"server_ip"`
	ServerPort int `json:"server_port"`
	ServerUser string `json:"server_user"`
	ServerPassword string `json:"server_password"`

	Host    string `json:"host"`
	Port    int    `json:"port"`
	UserId  string `json:"user_id"`
	Sni     string `json:"sni"`
	AlterId int    `json:"alter_id"`
	UseTls  bool   `json:"use_tls"`

	TransportType int `json:"transport_type"`
	Tcp           struct {
		Type string `json:"type"`
	} `json:"tcp"`
	WebSocket struct {
		Path string `json:"path"`
	} `json:"web_socket"`
	Kcp struct {
		Type string `json:"type"`
	} `json:"kcp"`
	Http2 struct {
		Host string `json:"host"`
		Path string `json:"path"`
	} `json:"http2"`
}

func V2rayServerDeploy(c *gin.Context) {}

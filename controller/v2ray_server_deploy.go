package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/Luna-CY/v2ray-helper/common/v2ray"
	"github.com/Luna-CY/v2ray-helper/common/webserver"
	"github.com/Luna-CY/v2ray-helper/dataservice"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"runtime"
	"strings"
	"time"
)

type V2rayServerDeployForm struct {
	ServerType  int `json:"server_type"`
	InstallType int `json:"install_type"`

	UseTls  bool   `json:"use_tls"`
	TlsHost string `json:"tls_host"`

	UseCloudreve bool `json:"use_cloudreve"`

	V2rayConfig v2ray.Config `json:"v2ray_config"`
}

const (
	ServerTypeLocalServer = iota + 1
	ServerTypeRemoteServer
)

const (
	InstallTypeDefault = iota + 1
	InstallTypeForce
	InstallTypeUpgrade
	InstallTypeReConfig
)

func V2rayServerDeploy(c *gin.Context) {
	var body V2rayServerDeployForm
	if err := c.ShouldBind(&body); nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	if !configurator.GetMainConfig().AllowV2rayDeploy {
		response.Response(c, code.BadRequest, "当前服务器已禁止部署V2ray", nil)

		return
	}

	if ServerTypeRemoteServer == body.ServerType {
		response.Response(c, code.BadRequest, "暂未支持远程服务器安装", nil)

		return
	}

	if !v2ray.CheckSystem(runtime.GOOS, runtime.GOARCH) {
		response.Response(c, code.BadRequest, fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH), nil)

		return
	}

	if err := v2ray.CheckExists(v2ray.CmdPath); nil != err {
		if !os.IsExist(err) {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "检查是否存在V2ray失败，详细请查看日志", nil)

			return
		}

		if InstallTypeDefault == body.InstallType {
			response.Response(c, code.BadRequest, "该服务器已安装V2ray，不能重复安装", nil)

			return
		}
	}

	body, err := v2rayServerDeployBodyFilter(body)
	if nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.BadRequest, "过滤输入失败，详细请查看日志", nil)

		return
	}

	// 如果有Nginx服务器并且已启动，那么停止Nginx，否则Caddy无法启动
	nginxIsRunning, err := webserver.CheckNginxIsRunning()
	if nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "检查Nginx服务状态失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

		return
	}

	if nginxIsRunning {
		if err := webserver.StopNginx(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "停止Nginx服务失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	v2rayIsRunning, err := v2ray.IsRunning()
	if nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "检查V2ray服务状态失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

		return
	}

	// 如果服务正在运行必须先停止V2ray服务，否则无法重新安装
	if v2rayIsRunning {
		if err := v2ray.Stop(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "停止V2ray服务失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	caddyIsRunning, err := webserver.CheckCaddyIsRunning()
	if nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "检查Caddy服务状态失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

		return
	}

	// 如果Caddy已启动需要停止服务，否则无法重新安装
	if caddyIsRunning {
		if err := webserver.StopCaddy(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "停止Caddy服务失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	// 仅在默认安装、强制安装、仅升级V2ray时安装V2ray
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType || InstallTypeUpgrade == body.InstallType {
		if err := v2ray.InstallLastRelease(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装V2ray失败，详细请查看日志", nil)

			return
		}
	}

	// 仅在默认安装、强制安装与仅配置V2ray时配置V2ray
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType || InstallTypeReConfig == body.InstallType {
		if err := v2ray.SetConfig(v2ray.ConfigPath, &body.V2rayConfig); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "配置V2ray失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	// 启动V2ray服务
	if err := v2ray.Start(); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "启动V2ray服务失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

		return
	}

	// 如果使用的传输类型是WebSocket，需要安装Caddy
	// 仅在默认安装与强制安装时配置Caddy
	if v2ray.TransportTypeWebSocket == body.V2rayConfig.TransportType && (InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType) {
		if err := webserver.InstallCaddyLastRelease(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装Caddy失败，详细请查看日志", nil)

			return
		}

		proxyPath := "/"
		if v2ray.TransportTypeWebSocket == body.V2rayConfig.TransportType {
			proxyPath = body.V2rayConfig.WebSocket.Path
		}
		if v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType {
			proxyPath = body.V2rayConfig.Http2.Path
		}

		if err := webserver.AppendCaddyConfigOnlyV2rayToSystem(body.TlsHost, body.UseTls, body.V2rayConfig.V2rayPort, proxyPath); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装Caddy失败，详细请查看日志", nil)

			return
		}
	}

	// 启动Caddy服务
	// WebSocket需要启动Caddy
	if v2ray.TransportTypeWebSocket == body.V2rayConfig.TransportType {
		if err := webserver.StartCaddy(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "启动Caddy服务失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	// 仅在默认安装与强制安装时处理Cloudreve
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType {
		// 仅在安装Cloudreve且V2ray传入类型为Websocket或HTTP2时配置Cloudreve
		if body.UseCloudreve && (v2ray.TransportTypeWebSocket == body.V2rayConfig.TransportType || v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType) {
			// TODO
		}
	}

	if err := generateConfig(body); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "生成客户端配置失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

		return
	}

	response.Success(c, code.OK, nil)
}

// v2rayServerDeployBodyFilter 表单过滤
func v2rayServerDeployBodyFilter(body V2rayServerDeployForm) (V2rayServerDeployForm, error) {
	for i, client := range body.V2rayConfig.Clients {
		if "" == client.UserId {
			id, err := uuid.NewRandom()
			if nil != err {
				return body, errors.New(fmt.Sprintf("无法生成用户ID: %v", err))
			}

			client.UserId = id.String()
		}

		body.V2rayConfig.Clients[i] = client
	}

	body.V2rayConfig.Tcp.Type = strings.TrimSpace(body.V2rayConfig.Tcp.Type)
	body.V2rayConfig.Tcp.Request.Version = strings.TrimSpace(body.V2rayConfig.Tcp.Request.Version)
	body.V2rayConfig.Tcp.Request.Method = strings.TrimSpace(body.V2rayConfig.Tcp.Request.Method)

	body.V2rayConfig.Tcp.Request.Path = strings.TrimSpace(body.V2rayConfig.Tcp.Request.Path)
	if "" == body.V2rayConfig.Tcp.Request.Path {
		body.V2rayConfig.Tcp.Request.Path = "/"
	}

	for i, h := range body.V2rayConfig.Tcp.Request.Headers {
		h.Key = strings.TrimSpace(h.Key)
		h.Value = strings.TrimSpace(h.Value)

		body.V2rayConfig.Tcp.Request.Headers[i] = h
	}

	body.V2rayConfig.Tcp.Response.Version = strings.TrimSpace(body.V2rayConfig.Tcp.Response.Version)
	body.V2rayConfig.Tcp.Response.Reason = strings.TrimSpace(body.V2rayConfig.Tcp.Response.Reason)

	for i, h := range body.V2rayConfig.Tcp.Response.Headers {
		h.Key = strings.TrimSpace(h.Key)
		h.Value = strings.TrimSpace(h.Value)

		if "" != h.Value {
			tokens := strings.Split(h.Value, ";;;")
			for i, s := range tokens {
				tokens[i] = strings.TrimSpace(s)
			}

			h.Value = strings.Join(tokens, ";;;")
		}

		body.V2rayConfig.Tcp.Response.Headers[i] = h
	}

	body.V2rayConfig.WebSocket.Path = strings.TrimSpace(body.V2rayConfig.WebSocket.Path)
	for i, h := range body.V2rayConfig.WebSocket.Headers {
		h.Key = strings.TrimSpace(h.Key)
		h.Value = strings.TrimSpace(h.Value)

		if "" != h.Value {
			tokens := strings.Split(h.Value, ";;;")
			for i, s := range tokens {
				tokens[i] = strings.TrimSpace(s)
			}

			h.Value = strings.Join(tokens, ";;;")
		}

		body.V2rayConfig.WebSocket.Headers[i] = h
	}

	body.V2rayConfig.Kcp.Type = strings.TrimSpace(body.V2rayConfig.Kcp.Type)

	body.V2rayConfig.Http2.Host = strings.TrimSpace(body.V2rayConfig.Http2.Host)
	if "" != body.V2rayConfig.Http2.Host {
		tokens := strings.Split(body.V2rayConfig.Http2.Host, ",")
		for i, s := range tokens {
			tokens[i] = strings.TrimSpace(s)
		}

		body.V2rayConfig.Http2.Host = strings.Join(tokens, ",")
	}

	body.V2rayConfig.Http2.Path = strings.TrimSpace(body.V2rayConfig.Http2.Path)
	if "" == body.V2rayConfig.Http2.Path {
		body.V2rayConfig.Http2.Path = "/"
	}

	body.TlsHost = strings.TrimSpace(body.TlsHost)

	return body, nil
}

// generateConfig 生成客户端配置
func generateConfig(body V2rayServerDeployForm) error {
	tcp, err := json.Marshal(body.V2rayConfig.Tcp)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	webSocket, err := json.Marshal(body.V2rayConfig.WebSocket)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	kcp, err := json.Marshal(body.V2rayConfig.Kcp)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	http2, err := json.Marshal(body.V2rayConfig.Http2)
	if nil != err {
		return errors.New(fmt.Sprintf("序列化数据失败: %v", err))
	}

	tcpString := string(tcp)
	webSocketString := string(webSocket)
	kcpString := string(kcp)
	http2String := string(http2)

	host := body.TlsHost
	if !body.UseTls {
		ip, err := util.GetPublicIpv4()
		if nil != err {
			return errors.New(fmt.Sprintf("获取本机IP失败: %v", err))
		}

		host = ip
	}

	port := 80
	if body.UseTls {
		port = 443
	}

	if v2ray.TransportTypeTcp == body.V2rayConfig.TransportType || v2ray.TransportTypeKcp == body.V2rayConfig.TransportType {
		port = body.V2rayConfig.V2rayPort
	}

	useTls := 0
	if body.UseTls {
		useTls = 1
	}

	one := 1

	for _, client := range body.V2rayConfig.Clients {
		endpoint := model.V2rayEndpoint{
			Cloud:         &one,
			Endpoint:      &one,
			Host:          &host,
			Port:          &port,
			UserId:        &client.UserId,
			AlterId:       &client.AlterId,
			UseTls:        &useTls,
			TransportType: &body.V2rayConfig.TransportType,
			Tcp:           &tcpString,
			WebSocket:     &webSocketString,
			Kcp:           &kcpString,
			Http2:         &http2String,
		}

		ct := time.Now().Unix()
		endpoint.CreateTime = &ct
		endpoint.Deleted = util.NewFalsePtr()

		if err := dataservice.GetBaseService().Create(&endpoint); nil != err {
			return errors.New(fmt.Sprintf("保存数据失败: %v", err))
		}
	}

	return nil
}

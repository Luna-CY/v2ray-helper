package controller

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/deploy"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/software/cloudreve"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"runtime"
	"strings"
)

type V2rayServerDeployForm struct {
	ServerType  int `json:"server_type"`
	InstallType int `json:"install_type"`

	UseTls  bool   `json:"use_tls"`
	TlsHost string `json:"tls_host"`

	EnableWebService bool   `json:"enable_web_service"`
	WebServiceType   string `json:"web_service_type"`

	CloudreveConfig struct {
		EnableAria2        bool `json:"enable_aria2"`
		ResetAdminPassword bool `json:"reset_admin_password"`
	} `json:"cloudreve_config"`

	V2rayConfig *v2ray.Config `json:"v2ray_config"`
}

const (
	ServerTypeLocalServer = iota + 1
	ServerTypeRemoteServer
)

const (
	InstallTypeDefault = iota + 1
	InstallTypeForce
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
			response.Response(c, code.BadRequest, "该服务器已安装V2ray，不能重复安装。请选择重新安装", nil)

			return
		}
	}

	body, err := v2rayServerDeployBodyFilter(body)
	if nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.BadRequest, "过滤输入失败，详细请查看日志", nil)

		return
	}

	if body.EnableWebService && (v2ray.TransportTypeTcp == body.V2rayConfig.TransportType || v2ray.TransportTypeKcp == body.V2rayConfig.TransportType) {
		response.Response(c, code.BadRequest, "使用TCP或KCP模式时不能开启站点伪装", nil)

		return
	}

	if body.EnableWebService && v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType && "/" == body.V2rayConfig.Http2.Path {
		response.Response(c, code.BadRequest, "开启站点伪装时，HTTP2的路径不能为/", nil)

		return
	}

	if v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType && !body.UseTls {
		response.Response(c, code.BadRequest, "使用HTTP2模式时必须开启HTTPS选项", nil)

		return
	}

	if body.UseTls && "" == body.TlsHost {
		response.Response(c, code.BadRequest, "开启HTTPS时必须填写域名，且域名必须绑定到此服务器", nil)

		return
	}

	config := &deploy.Config{}
	config.HttpsHost = body.TlsHost
	config.V2rayConfig = body.V2rayConfig

	if body.EnableWebService && cloudreve.Name == body.WebServiceType {
		config.FakeConfig = &deploy.FakeWebServerConfig{FakeType: deploy.FakeTypeCloudreve}
		config.FakeConfig.CloudreveConfig = &deploy.CloudreveConfig{EnableAria2: body.CloudreveConfig.EnableAria2, ResetAdminPassword: body.CloudreveConfig.ResetAdminPassword}
	}

	if err := deploy.GetManager().DeployServer(config); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.BadRequest, "启动部署失败，详细请查看日志", nil)

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

			client.UserId = strings.ToUpper(id.String())
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
	body.WebServiceType = strings.TrimSpace(body.WebServiceType)

	return body, nil
}

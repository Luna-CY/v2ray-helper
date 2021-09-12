package controller

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/v2ray"
	"github.com/gin-gonic/gin"
	"runtime"
	"strings"
)

type V2rayServerDeployForm struct {
	ServerType  int `json:"server_type"`
	InstallType int `json:"install_type"`

	UseTls  bool   `json:"use_tls"`
	TlsHost string `json:"tls_host"`

	V2rayConfig v2ray.Config `json:"v2ray_config"`
}

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

	if !v2ray.CheckSystem(runtime.GOOS, runtime.GOARCH) {
		response.Response(c, code.BadRequest, fmt.Sprintf("未受支持的系统: %v %v", runtime.GOOS, runtime.GOARCH), nil)

		return
	}

	exists, err := v2ray.CheckExists(v2ray.CmdPath)
	if nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "检查是否存在V2ray失败，详细请查看日志", nil)

		return
	}

	if exists && InstallTypeDefault == body.InstallType {
		response.Response(c, code.BadRequest, "该服务器已安装V2ray，不能重复安装", nil)

		return
	}

	body = v2rayServerDeployBodyFilter(body)

	// TODO 申请HTTPS证书

	if InstallTypeReConfig != body.InstallType {
		if err := v2ray.InstallLastRelease(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装V2ray失败，详细请查看日志", nil)

			return
		}
	}

	if err := v2ray.SetConfig(v2ray.ConfigPath, &body.V2rayConfig); nil != err {
		if err := v2ray.InstallLastRelease(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "配置V2ray失败，详细请查看日志。请使用强制安装或重新配置安装方式", nil)

			return
		}
	}

	response.Success(c, code.OK, nil)
}

// v2rayServerDeployBodyFilter 表单过滤
func v2rayServerDeployBodyFilter(body V2rayServerDeployForm) V2rayServerDeployForm {
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

		body.V2rayConfig.Tcp.Response.Headers[i] = h
	}

	body.V2rayConfig.WebSocket.Path = strings.TrimSpace(body.V2rayConfig.WebSocket.Path)
	for i, h := range body.V2rayConfig.WebSocket.Headers {
		h.Key = strings.TrimSpace(h.Key)
		h.Value = strings.TrimSpace(h.Value)

		body.V2rayConfig.WebSocket.Headers[i] = h
	}

	body.V2rayConfig.Kcp.Type = strings.TrimSpace(body.V2rayConfig.Kcp.Type)

	body.V2rayConfig.Http2.Host = strings.TrimSpace(body.V2rayConfig.Http2.Host)
	body.V2rayConfig.Http2.Path = strings.TrimSpace(body.V2rayConfig.Http2.Path)
	if "" == body.V2rayConfig.Http2.Path {
		body.V2rayConfig.Http2.Path = "/"
	}

	body.TlsHost = strings.TrimSpace(body.TlsHost)

	return body
}

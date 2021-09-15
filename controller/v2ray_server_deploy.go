package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/software/caddy"
	"github.com/Luna-CY/v2ray-helper/common/software/cloudreve"
	"github.com/Luna-CY/v2ray-helper/common/software/nginx"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
	"github.com/Luna-CY/v2ray-helper/common/util"
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

	EnableWebService bool   `json:"enable_web_service"`
	WebServiceType   string `json:"web_service_type"`

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

	if err := stopAllService(); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.BadRequest, "坚持Nginx/Caddy/V2ray/Cloudreve服务状态失败，详细请查看日志", nil)

		return
	}

	result := map[string]interface{}{}

	// WebSocket通过Caddy自动申请证书
	if body.UseTls && v2ray.TransportTypeWebSocket != body.V2rayConfig.TransportType {
		// 如果证书不存在先申请证书
		if !certificate.GetManager().CheckExists(body.TlsHost) {
			cert, err := certificate.GetManager().IssueNew(body.TlsHost, configurator.GetMainConfig().Email)
			if nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "申请HTTPS证书失败，详细请查看日志", nil)

				return
			}

			body.V2rayConfig.TlsKeyFile = cert.GetPrivateKeyFilePath()
			body.V2rayConfig.TlsCertFile = cert.GetCertificateFilePath()
		}

		body.V2rayConfig.UseTls = true
		body.V2rayConfig.TlsHost = body.TlsHost
	}

	// 仅在默认安装、重新安装、仅升级V2ray时安装V2ray
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType || InstallTypeUpgrade == body.InstallType {
		if err := v2ray.InstallLastRelease(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装V2ray失败，详细请查看日志", nil)

			return
		}
	}

	// 仅在默认安装、重新安装与仅配置V2ray时配置V2ray
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType || InstallTypeReConfig == body.InstallType {
		if err := v2ray.SetConfig(v2ray.ConfigPath, &body.V2rayConfig); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "配置V2ray失败，详细请查看日志。请使用重新安装的方式重试", nil)

			return
		}
	}

	// 启动V2ray服务
	if err := v2ray.Start(); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "启动V2ray服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

		return
	}

	if err := v2ray.Enable(); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "安装V2ray失败，详细请查看日志", nil)

		return
	}

	// 如果使用的传输类型是WebSocket，需要安装Caddy
	// 仅在默认安装与重新安装时配置Caddy
	if v2ray.TransportTypeWebSocket == body.V2rayConfig.TransportType || v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType {
		if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType {
			if err := caddy.InstallLastRelease(); nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "安装Caddy失败，详细请查看日志", nil)

				return
			}

			port := caddy.PortHttp
			if body.UseTls {
				port = caddy.PortHttps
			}

			enableCloudreve := false
			if body.EnableWebService {
				enableCloudreve = true
			}

			enableHttp2 := false
			if v2ray.TransportTypeHttp2 == body.V2rayConfig.TransportType {
				enableHttp2 = true
			}

			if err := caddy.SetConfigToSystem(body.TlsHost, port, body.V2rayConfig.V2rayPort, body.V2rayConfig.WebSocket.Path, enableCloudreve, enableHttp2); nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "安装Caddy失败，详细请查看日志", nil)

				return
			}
		}

		// 启动Caddy服务
		if err := caddy.Start(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "启动Caddy服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

			return
		}

		if err := caddy.Enable(); nil != err {
			logger.GetLogger().Errorln(err)

			response.Response(c, code.ServerError, "安装Caddy失败，详细请查看日志", nil)

			return
		}
	}

	// 处理站点伪装配置
	if InstallTypeDefault == body.InstallType || InstallTypeForce == body.InstallType {
		if body.EnableWebService && cloudreve.Name == body.WebServiceType {
			if err := cloudreve.InstallLastRelease(); nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "安装Cloudreve服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

				return
			}

			if err := cloudreve.Enable(); nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "安装Cloudreve服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

				return
			}

			if err := cloudreve.Start(); nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "启动Cloudreve服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

				return
			}

			password, err := cloudreve.ResetAdminPassword()
			if nil != err {
				logger.GetLogger().Errorln(err)

				response.Response(c, code.ServerError, "安装Cloudreve服务失败，详细请查看日志。请使用重新安装的方式重试", nil)

				return
			}

			result["cloudreve_admin"] = "admin@cloudreve.org"
			result["cloudreve_password"] = password
		}
	}

	if err := generateConfig(body); nil != err {
		logger.GetLogger().Errorln(err)

		response.Response(c, code.ServerError, "生成客户端配置失败，详细请查看日志。请使用重新安装的方式重试", nil)

		return
	}

	data := gin.H(result)
	response.Success(c, code.OK, &data)
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

// stopAllService 停止所有服务
func stopAllService() error {
	// 如果有Nginx服务器并且已启动，那么停止Nginx，否则Caddy无法启动
	nginxIsRunning, err := nginx.IsRunning()
	if nil != err {
		return err
	}

	if nginxIsRunning {
		if err := nginx.Stop(); nil != err {
			return err
		}
	}

	caddyIsRunning, err := caddy.IsRunning()
	if nil != err {
		return err
	}

	// 如果Caddy已启动需要停止服务，否则无法重新安装
	if caddyIsRunning {
		if err := caddy.Stop(); nil != err {
			return err
		}
	}

	v2rayIsRunning, err := v2ray.IsRunning()
	if nil != err {
		return err
	}

	// 如果服务正在运行必须先停止V2ray服务，否则无法重新安装
	if v2rayIsRunning {
		if err := v2ray.Stop(); nil != err {
			return err
		}
	}

	cloudreveIsRunning, err := cloudreve.IsRunning()
	if nil != err {
		return err
	}

	if cloudreveIsRunning {
		if err := cloudreve.Stop(); nil != err {
			return err
		}
	}

	return nil
}

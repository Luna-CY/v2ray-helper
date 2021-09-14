package controller

import (
	"encoding/json"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/software/v2ray"
	"github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/Luna-CY/v2ray-helper/dataservice"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type V2rayEndpointNewForm struct {
	Remark string `json:"remark"`

	Host    string `json:"host"`
	Port    int    `json:"port"`
	UserId  string `json:"user_id"`
	AlterId int    `json:"alter_id"`
	UseTls  bool   `json:"use_tls"`

	TransportType int `json:"transport_type"`
	Tcp           struct {
		Type    string `json:"type"`
		Request struct {
			Version string               `json:"version"`
			Method  string               `json:"method"`
			Path    string               `json:"path"`
			Headers []v2ray.ConfigHeader `json:"headers"`
		} `json:"request"`
		Response struct {
			Version string               `json:"version"`
			Status  string               `json:"status"`
			Reason  string               `json:"reason"`
			Headers []v2ray.ConfigHeader `json:"headers"`
		} `json:"response"`
	} `json:"tcp"`
	WebSocket struct {
		Path    string               `json:"path"`
		Headers []v2ray.ConfigHeader `json:"headers"`
	} `json:"web_socket"`
	Kcp struct {
		Type             string `json:"type"`
		Mtu              int    `json:"mtu"`
		Tti              int    `json:"tti"`
		UpLinkCapacity   int    `json:"uplink_capacity"`
		DownLinkCapacity int    `json:"downlink_capacity"`
		Congestion       bool   `json:"congestion"`
		ReadBufferSize   int    `json:"read_buffer_size"`
		WriteBufferSize  int    `json:"write_buffer_size"`
	} `json:"kcp"`
	Http2 struct {
		Host string `json:"host"`
		Path string `json:"path"`
	} `json:"http2"`
}

func V2rayEndpointNew(c *gin.Context) {
	var body V2rayEndpointNewForm
	if err := c.ShouldBind(&body); nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	body.Remark = strings.TrimSpace(body.Remark)
	body.Host = strings.TrimSpace(body.Host)
	body.UserId = strings.TrimSpace(body.UserId)
	body.Tcp.Type = strings.TrimSpace(body.Tcp.Type)
	body.WebSocket.Path = strings.TrimSpace(body.WebSocket.Path)
	body.Kcp.Type = strings.TrimSpace(body.Kcp.Type)
	body.Http2.Host = strings.TrimSpace(body.Http2.Host)
	body.Http2.Path = strings.TrimSpace(body.Http2.Path)

	if "" == body.Host || 0 == body.Port || "" == body.UserId || 0 == body.TransportType {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	tcp, err := json.Marshal(body.Tcp)
	if nil != err {
		logger.GetLogger().Errorf("序列化数据失败: %v", err)

		response.Response(c, code.ServerError, "序列化数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	webSocket, err := json.Marshal(body.WebSocket)
	if nil != err {
		logger.GetLogger().Errorf("序列化数据失败: %v", err)

		response.Response(c, code.ServerError, "序列化数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	kcp, err := json.Marshal(body.Kcp)
	if nil != err {
		logger.GetLogger().Errorf("序列化数据失败: %v", err)

		response.Response(c, code.ServerError, "序列化数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	http2, err := json.Marshal(body.Http2)
	if nil != err {
		logger.GetLogger().Errorf("序列化数据失败: %v", err)

		response.Response(c, code.ServerError, "序列化数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	useTls := 0
	if body.UseTls {
		useTls = 1
	}

	tcpString := string(tcp)
	webSocketString := string(webSocket)
	kcpString := string(kcp)
	http2String := string(http2)

	one := 1

	endpoint := model.V2rayEndpoint{
		Cloud:         &one,
		Endpoint:      &one,
		Remark:        &body.Remark,
		Host:          &body.Host,
		Port:          &body.Port,
		UserId:        &body.UserId,
		AlterId:       &body.AlterId,
		UseTls:        &useTls,
		TransportType: &body.TransportType,
		Tcp:           &tcpString,
		WebSocket:     &webSocketString,
		Kcp:           &kcpString,
		Http2:         &http2String,
	}

	ct := time.Now().Unix()
	endpoint.CreateTime = &ct
	endpoint.Deleted = util.NewFalsePtr()

	if err := dataservice.GetBaseService().Create(&endpoint); nil != err {
		logger.GetLogger().Errorf("保存数据失败: %v", err)

		response.Response(c, code.ServerError, "保存数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	response.Success(c, code.OK, nil)
}

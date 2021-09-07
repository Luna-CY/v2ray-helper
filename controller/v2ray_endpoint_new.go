package controller

import (
	"encoding/json"
	"gitee.com/Luna-CY/v2ray-subscription/code"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"gitee.com/Luna-CY/v2ray-subscription/dataservice"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/response"
	"gitee.com/Luna-CY/v2ray-subscription/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type V2rayEndpointNewForm struct {
	Cloud    int    `json:"cloud"`
	Endpoint int    `json:"endpoint"`
	Rate     string `json:"rate"`
	Remark   string `json:"remark"`

	Host          string `json:"host"`
	Port          int    `json:"port"`
	UserId        string `json:"user_id"`
	AlterId       int    `json:"alter_id"`
	Level         int    `json:"level"`
	TransportType int    `json:"transport_type"`
	WebSocket     struct {
		Path string `json:"path"`
	} `json:"web_socket"`
}

func V2rayEndpointNew(c *gin.Context) {
	var body V2rayEndpointNewForm
	if err := c.ShouldBind(&body); nil != err {
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	body.Rate = strings.TrimSpace(body.Rate)
	body.Remark = strings.TrimSpace(body.Remark)
	body.Host = strings.TrimSpace(body.Host)
	body.UserId = strings.TrimSpace(body.UserId)
	body.WebSocket.Path = strings.TrimSpace(body.WebSocket.Path)

	if 0 == body.Cloud || 0 == body.Endpoint || "" == body.Host || 0 == body.Port || "" == body.UserId || 0 == body.TransportType {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	webSocket, err := json.Marshal(body.WebSocket)
	if nil != err {
		logger.GetLogger().Errorf("序列化数据失败: %v", err)

		response.Response(c, code.ServerError, "序列化数据失败，稍后重试一下吧，或者联系管理员", nil)

		return
	}

	webSocketString := string(webSocket)
	endpoint := model.V2rayEndpoint{
		Cloud:         &body.Cloud,
		Endpoint:      &body.Endpoint,
		Rate:          &body.Rate,
		Remark:        &body.Remark,
		Host:          &body.Host,
		Port:          &body.Port,
		UserId:        &body.UserId,
		AlterId:       &body.AlterId,
		Level:         &body.Level,
		TransportType: &body.TransportType,
		WebSocket:     &webSocketString,
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

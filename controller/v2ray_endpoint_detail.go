package controller

import (
	"encoding/json"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/dataservice"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type V2rayEndpointDetailQuery struct {
	Id uint `form:"id"`
}

func V2rayEndpointDetail(c *gin.Context) {
	var query V2rayEndpointDetailQuery
	if err := c.ShouldBindQuery(&query); nil != err {
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	if 0 == query.Id {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	var endpoint model.V2rayEndpoint
	if err := dataservice.GetBaseService().GetById(query.Id, &endpoint); nil != err {
		if gorm.ErrRecordNotFound != err {
			logger.GetLogger().Errorf("查询数据库失败: %v", err)
		}

		response.Response(c, code.BadRequest, "查询数据失败", nil)

		return
	}

	res := map[string]interface{}{
		"host":           *endpoint.Host,
		"port":           *endpoint.Port,
		"user_id":        *endpoint.UserId,
		"alter_id":       *endpoint.AlterId,
		"use_tls":        false,
		"transport_type": *endpoint.TransportType,
		"tcp":            gin.H{},
		"web_socket":     gin.H{},
		"kcp":            gin.H{},
		"http2":          gin.H{},
	}

	if 1 == *endpoint.UseTls {
		res["use_tls"] = true
	}

	if model.V2rayEndpointTransportTypeTcp == *endpoint.TransportType {
		var tcp model.V2rayEndpointTcp
		if err := json.Unmarshal([]byte(*endpoint.Tcp), &tcp); nil != err {
			response.Response(c, code.BadRequest, "无法解析此配置，请稍后重试一下吧，或联系管理员", nil)

			return
		}

		res["tcp"] = tcp
	}

	if model.V2rayEndpointTransportTypeWebSocket == *endpoint.TransportType {
		var webSocket model.V2rayEndpointWebSocket
		if err := json.Unmarshal([]byte(*endpoint.WebSocket), &webSocket); nil != err {
			response.Response(c, code.BadRequest, "无法解析此配置，请稍后重试一下吧，或联系管理员", nil)

			return
		}

		res["web_socket"] = webSocket
	}

	if model.V2rayEndpointTransportTypeKcp == *endpoint.TransportType {
		var kcp model.V2rayEndpointKcp
		if err := json.Unmarshal([]byte(*endpoint.Kcp), &kcp); nil != err {
			response.Response(c, code.BadRequest, "无法解析此配置，请稍后重试一下吧，或联系管理员", nil)

			return
		}

		res["kcp"] = kcp
	}

	if model.V2rayEndpointTransportTypeHttp2 == *endpoint.TransportType {
		var http2 model.V2rayEndpointHttp2
		if err := json.Unmarshal([]byte(*endpoint.Http2), &http2); nil != err {
			response.Response(c, code.BadRequest, "无法解析此配置，请稍后重试一下吧，或联系管理员", nil)

			return
		}

		res["http2"] = http2
	}

	data := gin.H(res)
	response.Success(c, code.OK, &data)
}

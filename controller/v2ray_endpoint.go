package controller

import (
	"encoding/json"
	"gitee.com/Luna-CY/v2ray-subscription/code"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"gitee.com/Luna-CY/v2ray-subscription/dataservice"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/response"
	"github.com/gin-gonic/gin"
)

func V2rayEndpointList(c *gin.Context) {
	var v2rayEndpointList []model.V2rayEndpoint
	if err := dataservice.GetBaseService().FindByCondition(&v2rayEndpointList, "id desc", "deleted = 0"); nil != err {
		response.Response(c, code.ServerError, "查询数据失败", nil)

		return
	}

	result := response.NewEmptyJsonList()
	for _, endpoint := range v2rayEndpointList {
		var ws gin.H
		if nil != endpoint.WebSocket && "" != *endpoint.WebSocket {
			if err := json.Unmarshal([]byte(*endpoint.WebSocket), &ws); nil != err {
				logger.GetLogger().Errorf("无法解析WebSocket配置: %v", err)

				continue
			}
		}

		result = append(result, gin.H{
			"id":             endpoint.Id,
			"cloud":          endpoint.Cloud,
			"endpoint":       endpoint.Endpoint,
			"rate":           endpoint.Rate,
			"remark":         endpoint.Remark,
			"host":           endpoint.Host,
			"transport_type": endpoint.TransportType,
		})
	}

	response.Success(c, code.OK, &gin.H{"data": result})
}

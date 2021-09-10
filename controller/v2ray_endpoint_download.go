package controller

import (
	"gitee.com/Luna-CY/v2ray-subscription/code"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"gitee.com/Luna-CY/v2ray-subscription/dataservice"
	"gitee.com/Luna-CY/v2ray-subscription/generator"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type V2rayEndpointDownloadBody struct {
	Id   uint `json:"id"`
	Type int  `json:"type"` // 下载类型：1表示V2rayX；2表示V2rayNG
}

func V2rayEndpointDownload(c *gin.Context) {
	var body V2rayEndpointDownloadBody
	if err := c.ShouldBind(&body); nil != err {
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	if 0 == body.Id {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	var endpoint model.V2rayEndpoint
	if err := dataservice.GetBaseService().GetById(body.Id, &endpoint); nil != err {
		if gorm.ErrRecordNotFound != err {
			logger.GetLogger().Errorf("查询数据库失败: %v", err)
		}

		response.Response(c, code.BadRequest, "查询数据失败", nil)

		return
	}

	content, err := generator.GenerateVMessShareLink(endpoint)
	if nil != err {
		response.Response(c, code.ServerError, err.Error(), nil)

		return
	}

	response.Success(c, code.OK, &gin.H{"content": content})
}
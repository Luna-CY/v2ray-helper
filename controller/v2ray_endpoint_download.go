package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/generator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/gin-gonic/gin"
)

type V2rayEndpointDownloadBody struct {
	Id uint `json:"id"`
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

	content, err := generator.GenerateVMessShareLink()
	if nil != err {
		response.Response(c, code.ServerError, err.Error(), nil)

		return
	}

	response.Success(c, code.OK, &gin.H{"content": content})
}

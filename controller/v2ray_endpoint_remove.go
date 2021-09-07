package controller

import (
	"gitee.com/Luna-CY/v2ray-subscription/code"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"gitee.com/Luna-CY/v2ray-subscription/dataservice"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/response"
	"gitee.com/Luna-CY/v2ray-subscription/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type V2rayEndpointRemoveBody struct {
	Id       uint   `json:"id"`
	Password string `json:"password"`
}

func V2rayEndpointRemove(c *gin.Context) {
	var body V2rayEndpointRemoveBody
	if err := c.ShouldBind(&body); nil != err {
		response.Response(c, code.BadRequest, "错误请求", nil)

		return
	}

	body.Password = strings.TrimSpace(body.Password)

	if 0 == body.Id || "" == body.Password {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	if util.Md5(time.Now().Format("200601021504")) != body.Password {
		response.Response(c, code.BadRequest, "密码错误", nil)

		return
	}

	var endpoint model.V2rayEndpoint
	if err := dataservice.GetBaseService().GetById(body.Id, &endpoint); nil != err {
		if gorm.ErrRecordNotFound != err {
			logger.GetLogger().Errorf("查询数据库失败: %v", err)
		}

		response.Response(c, code.BadRequest, "找不到该节点信息", nil)

		return
	}

	endpoint.Deleted = util.NewTruePtr()

	if err := dataservice.GetBaseService().UpdateById(body.Id, &endpoint); nil != err {
		response.Response(c, code.ServerError, "删除失败，稍后重试一下吧，或联系管理员", nil)

		return
	}

	response.Success(c, code.OK, nil)
}

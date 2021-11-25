package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database/model"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	util2 "github.com/Luna-CY/v2ray-helper/common/util"
	"github.com/Luna-CY/v2ray-helper/dataservice"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
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

	if util2.Md5(viper.GetString(configurator.KeyAuthManagementKey)) != body.Password {
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

	endpoint.Deleted = util2.NewTruePtr()

	if err := dataservice.GetBaseService().UpdateById(body.Id, &endpoint); nil != err {
		response.Response(c, code.ServerError, "删除失败，稍后重试一下吧，或联系管理员", nil)

		return
	}

	response.Success(c, code.OK, nil)
}

package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/notice"
	"github.com/gin-gonic/gin"
)

func MetaInfo(c *gin.Context) {
	isDefaultAccessKey := false
	if configurator.DefaultAccessKey == configurator.GetMainConfig().AccessKey {
		isDefaultAccessKey = true
	}

	isDefaultManagementKey := false
	if configurator.DefaultManagementKey == configurator.GetMainConfig().ManagementKey {
		isDefaultManagementKey = true
	}

	data := &gin.H{
		"is_default_access_key":     isDefaultAccessKey,
		"is_default_management_key": isDefaultManagementKey,
		"listen":                    configurator.GetMainConfig().Listen,
		"enable_https":              configurator.GetMainConfig().EnableHttps,
		"https_host":                configurator.GetMainConfig().HttpsHost,
		"email":                     configurator.GetMainConfig().Email,
		"notice_list":               notice.GetManager().GetAll(),
	}
	response.Success(c, code.OK, data)
}

package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/notice"
	"github.com/gin-gonic/gin"
)

func MetaInfo(c *gin.Context) {
	isDefaultKey := false
	if configurator.DefaultKey == configurator.GetMainConfig().Key {
		isDefaultKey = true
	}

	isDefaultRemoveKey := false
	if configurator.DefaultRemoveKey == configurator.GetMainConfig().RemoveKey {
		isDefaultRemoveKey = true
	}

	data := &gin.H{
		"is_default_key":        isDefaultKey,
		"is_default_remove_key": isDefaultRemoveKey,
		"listen":                configurator.GetMainConfig().Listen,
		"enable_https":          configurator.GetMainConfig().EnableHttps,
		"https_host":            configurator.GetMainConfig().HttpsHost,
		"email":                 configurator.GetMainConfig().Email,
		"notice_list":           notice.GetManager().GetAll(),
	}
	response.Success(c, code.OK, data)
}

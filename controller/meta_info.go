package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/notice"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func MetaInfo(c *gin.Context) {
	isDefaultAccessKey := false
	if configurator.DefaultAccessKey == viper.GetString(configurator.KeyAuthAccessKey) {
		isDefaultAccessKey = true
	}

	isDefaultManagementKey := false
	if configurator.DefaultManagementKey == viper.GetString(configurator.KeyAuthManagementKey) {
		isDefaultManagementKey = true
	}

	data := &gin.H{
		"is_default_access_key":     isDefaultAccessKey,
		"is_default_management_key": isDefaultManagementKey,
		"listen":                    viper.GetInt(configurator.KeyServerPort),
		"enable_https":              viper.GetBool(configurator.KeyServerHttpsEnable),
		"https_host":                viper.GetString(configurator.KeyServerHttpsHost),
		"email":                     viper.GetString(configurator.KeyHttpsIssueEmail),
		"notice_list":               notice.GetManager().GetAll(),
	}

	response.Success(c, code.OK, data)
}

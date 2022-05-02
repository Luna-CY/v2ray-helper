package controller

import (
	"github.com/Luna-CY/v2ray-helper/v1/common/http/code"
	"github.com/Luna-CY/v2ray-helper/v1/common/http/response"
	"github.com/Luna-CY/v2ray-helper/v1/common/notice"
	"github.com/gin-gonic/gin"
)

func CleanNotice(c *gin.Context) {
	notice.GetManager().Clean()

	response.Success(c, code.OK, nil)
}

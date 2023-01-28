package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/gin-gonic/gin"
)

func CleanNotice(c *gin.Context) {
	response.Success(c, code.OK, nil)
}

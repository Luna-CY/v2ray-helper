package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success 响应一个成功消息
func Success(context *gin.Context, code int, data *gin.H) {
	Response(context, code, "OK", data)
}

// Response 响应一个错误消息
func Response(context *gin.Context, code int, message string, data *gin.H) {
	result := gin.H{
		"code":    code,
		"message": message,
	}

	if nil != data {
		result["data"] = data
	}

	context.JSON(http.StatusOK, result)
}

// NewEmptyJsonList 获取空的json列表
func NewEmptyJsonList() []gin.H {
	return make([]gin.H, 0, 0)
}

// NewDataListResult 统一数据列表响应结构
func NewDataListResult(data []gin.H, total int64, page, size int) *gin.H {
	return &gin.H{
		"data":       data,
		"pagination": Pagination(total, page, size),
	}
}

// Pagination 统一分页
func Pagination(total int64, page, size int) *gin.H {
	return &gin.H{
		"total": total,
		"page":  page,
		"size":  size,
	}
}

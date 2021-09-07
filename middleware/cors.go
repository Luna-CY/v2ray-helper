package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 同源跨域配置
func Cors(context *gin.Context) {
	context.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

	if http.MethodOptions == context.Request.Method {
		context.AbortWithStatus(http.StatusNoContent)

		return
	}

	context.Next()
}

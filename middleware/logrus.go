package middleware

import (
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LogRus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		logger.GetLogger().WithFields(logrus.Fields{
			"code":         statusCode,
			"process_time": latencyTime.Milliseconds(),
			"ip":           clientIP,
			"method":       reqMethod,
			"path":         reqUri,
		}).Info()
	}
}

package main

import (
	"fmt"
	"gitee.com/Luna-CY/v2ray-subscription/configurator"
	"gitee.com/Luna-CY/v2ray-subscription/database"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/middleware"
	"gitee.com/Luna-CY/v2ray-subscription/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	if err := configurator.Init(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	if err := logger.Init(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化日志失败: %v", err))
	}

	if err := database.Init(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	engine := gin.New()
	engine.Use(middleware.Cors).Use(middleware.LogRus())
	engine.LoadHTMLFiles("templates/index.html")

	engine.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{})
	})

	if err := router.RegisterApiRouter(engine.Group("/api")); nil != err {
		log.Fatalln(fmt.Sprintf("注册路由失败: %v", err))
	}

	if err := engine.Run("127.0.0.1:8800"); nil != err {
		log.Fatalf("启动服务器失败: %v\n", err)
	}
}

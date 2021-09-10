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
	engine.Use(middleware.LogRus())

	// 只有非Release模式才允许同源跨域
	if gin.ReleaseMode != gin.Mode() {
		engine.Use(middleware.Cors)
	}

	if err := router.RegisterApiRouter(engine.Group("/api")); nil != err {
		log.Fatalln(fmt.Sprintf("注册路由失败: %v", err))
	}

	if err := engine.Run(fmt.Sprintf("127.0.0.1:%v", configurator.GetMainConfig().Listen)); nil != err {
		log.Fatalf("启动服务器失败: %v\n", err)
	}
}

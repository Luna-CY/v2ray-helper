package router

import (
	"github.com/Luna-CY/v2ray-helper/controller"
	"github.com/Luna-CY/v2ray-helper/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(engine *gin.RouterGroup) error {
	jwt, err := middleware.GetJWT()
	if nil != err {
		return err
	}

	engine.POST("/auth", jwt.LoginHandler)

	engine.Use(jwt.MiddlewareFunc())
	engine.POST("/auth/logout", jwt.LogoutHandler)

	engine.GET("/v2ray/get-config", controller.CleanNotice)
	engine.POST("/v2ray/add-client", controller.CleanNotice)
	engine.POST("/v2ray/remove-client", controller.CleanNotice)
	engine.POST("/v2ray/add-inbound", controller.CleanNotice)
	engine.POST("/v2ray/remove-inbound", controller.CleanNotice)
	engine.POST("/v2ray/add-outbound", controller.CleanNotice)
	engine.POST("/v2ray/remove-outbound", controller.CleanNotice)
	engine.POST("/v2ray/upgrade", controller.CleanNotice)

	engine.GET("/cert/list", controller.CleanNotice)
	engine.POST("/cert/upload", controller.CleanNotice)
	engine.POST("/cert/issue", controller.CleanNotice)
	engine.POST("/cert/renew", controller.CleanNotice)
	engine.POST("/cert/remove", controller.CleanNotice)

	engine.GET("/host/list", controller.CleanNotice)
	engine.POST("/host/add", controller.CleanNotice)
	engine.POST("/host/update", controller.CleanNotice)
	engine.POST("/host/remove", controller.CleanNotice)

	engine.GET("/app/list", controller.CleanNotice)
	engine.POST("/app/install", controller.CleanNotice)
	engine.POST("/app/upgrade", controller.CleanNotice)
	engine.POST("/app/remove", controller.CleanNotice)

	return nil
}

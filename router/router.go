package router

import (
	"gitee.com/Luna-CY/v2ray-subscription/controller"
	"gitee.com/Luna-CY/v2ray-subscription/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(engine *gin.RouterGroup) error {
	jwt, err := middleware.GetJWT()
	if nil != err {
		return err
	}

	engine.POST("/auth", jwt.LoginHandler)
	engine.POST("/auth/logout", jwt.LogoutHandler)

	engine.Use(jwt.MiddlewareFunc())
	engine.GET("/v2ray-endpoint", controller.V2rayEndpointList)
	engine.POST("/v2ray-endpoint/new", controller.V2rayEndpointNew)
	engine.POST("/v2ray-endpoint/remove", controller.V2rayEndpointRemove)
	engine.POST("/v2ray-endpoint/download", controller.V2rayEndpointDownload)

	return nil
}

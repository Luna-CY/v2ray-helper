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

	engine.GET("/v2ray-endpoint", controller.V2rayEndpointList)
	engine.GET("/v2ray-endpoint/detail", controller.V2rayEndpointDetail)
	engine.POST("/v2ray-endpoint/new", controller.V2rayEndpointNew)
	engine.POST("/v2ray-endpoint/remove", controller.V2rayEndpointRemove)
	engine.POST("/v2ray-endpoint/download", controller.V2rayEndpointDownload)

	engine.POST("/v2ray-server-deploy", controller.V2rayServerDeploy)

	return nil
}

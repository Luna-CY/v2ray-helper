package router

import (
	controller2 "github.com/Luna-CY/v2ray-helper/v1/controller"
	"github.com/Luna-CY/v2ray-helper/v1/middleware"
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

	engine.GET("/meta-info", controller2.MetaInfo)
	engine.POST("/save-meta-info", controller2.SaveMetaInfo)

	engine.POST("/clean-notice", controller2.CleanNotice)

	engine.GET("/v2ray-endpoint", controller2.V2rayEndpointList)
	engine.GET("/v2ray-endpoint/detail", controller2.V2rayEndpointDetail)
	engine.POST("/v2ray-endpoint/new", controller2.V2rayEndpointNew)
	engine.POST("/v2ray-endpoint/remove", controller2.V2rayEndpointRemove)
	engine.POST("/v2ray-endpoint/download", controller2.V2rayEndpointDownload)

	engine.POST("/v2ray-server-deploy", controller2.V2rayServerDeploy)

	return nil
}

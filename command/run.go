package command

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/middleware"
	"github.com/Luna-CY/v2ray-helper/router"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/webjs"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "启动服务",
	Args:  cobra.NoArgs,
	Run:   run,
}

func run(*cobra.Command, []string) {
	// 自动续期证书
	go certificate.GetManager().RenewLoop()

	engine := gin.New()
	engine.Use(middleware.LogRus())

	// 只有非Release模式才允许同源跨域
	if gin.ReleaseMode != gin.Mode() {
		engine.Use(middleware.Cors)
	}

	if err := router.RegisterApiRouter(engine.Group("/api")); nil != err {
		log.Fatalln(fmt.Sprintf("注册路由失败: %v", err))
	}

	rootFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web"}
	cssFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web/css"}
	fontFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web/fonts"}

	javascriptFileSystem := assetfs.AssetFS{Asset: webjs.Asset, AssetDir: webjs.AssetDir, AssetInfo: webjs.AssetInfo, Prefix: "web/js"}

	engine.StaticFS("/js", &javascriptFileSystem)
	engine.StaticFS("/css", &cssFileSystem)
	engine.StaticFS("/fonts", &fontFileSystem)
	engine.StaticFS("/favicon.ico", &rootFileSystem)

	engine.GET("/manifest.json", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		c.Writer.Header().Set("Content-Type", "application/json")
		_, _ = c.Writer.Write(webstatic.MustAsset("web/manifest.json"))
		c.Writer.Flush()
	})

	engine.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		_, _ = c.Writer.Write(webstatic.MustAsset("web/index.html"))
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})

	address := fmt.Sprintf("%v:%v", viper.GetString(configurator.KeyServerAddress), viper.GetInt(configurator.KeyServerPort))
	if viper.GetBool(configurator.KeyServerHttpsEnable) {
		https, err := certificate.GetManager().GetCertificate(viper.GetString(configurator.KeyServerHttpsHost))
		if nil != err {
			log.Fatalln("无法获取HTTPS证书")
		}

		if err := engine.RunTLS(address, https.GetCertificateFilePath(), https.GetPrivateKeyFilePath()); nil != err {
			log.Fatalf("启动服务器失败: %v\n", err)
		}
	} else {
		if err := engine.Run(address); nil != err {
			log.Fatalf("启动服务器失败: %v\n", err)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/Luna-CY/v2ray-helper/middleware"
	"github.com/Luna-CY/v2ray-helper/router"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/img/imgclient"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/img/imghelp"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/img/imgicons"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/webjs"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	homeDir := ""

	flag.StringVar(&homeDir, "home-dir", "", "主目录，数据库及配置文件将放在此目录下，未指定时为当前目录")
	flag.Parse()

	homeDir = strings.TrimSpace(homeDir)

	if err := runtime.InitRuntime(runtime.GetRootPath(homeDir)); nil != err {
		log.Fatalln(fmt.Sprintf("初始化运行环境失败: %v", err))
	}

	if err := configurator.Init(runtime.GetRootPath(homeDir)); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	if err := logger.Init(runtime.GetRootPath(homeDir)); nil != err {
		log.Fatalln(fmt.Sprintf("初始化日志失败: %v", err))
	}

	if err := database.Init(filepath.Join(runtime.GetRootPath(homeDir), "main.db"), 10); nil != err {
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

	rootFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web"}
	cssFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web/css"}
	fontFileSystem := assetfs.AssetFS{Asset: webstatic.Asset, AssetDir: webstatic.AssetDir, AssetInfo: webstatic.AssetInfo, Prefix: "web/fonts"}

	javascriptFileSystem := assetfs.AssetFS{Asset: webjs.Asset, AssetDir: webjs.AssetDir, AssetInfo: webjs.AssetInfo, Prefix: "web/js"}

	imgClientFileSystem := assetfs.AssetFS{Asset: imgclient.Asset, AssetDir: imgclient.AssetDir, AssetInfo: imgclient.AssetInfo, Prefix: "web/img/client"}
	imgHelpFileSystem := assetfs.AssetFS{Asset: imghelp.Asset, AssetDir: imghelp.AssetDir, AssetInfo: imghelp.AssetInfo, Prefix: "web/img/help"}
	imgIconsFileSystem := assetfs.AssetFS{Asset: imgicons.Asset, AssetDir: imgicons.AssetDir, AssetInfo: imgicons.AssetInfo, Prefix: "web/img/icons"}

	engine.StaticFS("/js", &javascriptFileSystem)
	engine.StaticFS("/css", &cssFileSystem)
	engine.StaticFS("/fonts", &fontFileSystem)
	engine.StaticFS("/img/client", &imgClientFileSystem)
	engine.StaticFS("/img/help", &imgHelpFileSystem)
	engine.StaticFS("/img/icons", &imgIconsFileSystem)
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

	if err := engine.Run(configurator.GetMainConfig().GetListen()); nil != err {
		log.Fatalf("启动服务器失败: %v\n", err)
	}
}

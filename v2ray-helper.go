package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/database"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/Luna-CY/v2ray-helper/middleware"
	"github.com/Luna-CY/v2ray-helper/router"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic"
	"github.com/Luna-CY/v2ray-helper/staticfile/webstatic/webjs"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	homeDir := ""
	install := false
	installAndEnable := false

	flag.StringVar(&homeDir, "home-dir", "", "主目录，数据库及配置文件将放在此目录下，未指定时为当前目录")
	flag.BoolVar(&install, "install", false, "安装为系统服务并退出")
	flag.BoolVar(&installAndEnable, "install-and-enable", false, "安装为系统服务并且开机启动")
	flag.Parse()

	homeDir = strings.TrimSpace(homeDir)
	rootAbsPath := runtime.AbsRootPath(homeDir)

	if err := runtime.InitRuntime(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化运行环境失败: %v", err))
	}

	if err := configurator.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	// 设置GIN模式需要在其他组件初始化之前
	if configurator.GetMainConfig().GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// logger组件需要在其他组件之前初始化
	if err := logger.Init(rootAbsPath); nil != err {
		log.Fatalln(fmt.Sprintf("初始化日志失败: %v", err))
	}

	if err := certificate.Init(context.Background()); nil != err {
		log.Fatalln(fmt.Sprintf("初始化证书管理器失败: %v", err))
	}

	if err := database.Init(filepath.Join(rootAbsPath, "main.db"), 10); nil != err {
		log.Fatalln(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 安装为系统服务并退出
	if install || installAndEnable {
		installAsService(installAndEnable)

		os.Exit(0)
	}

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

	if err := engine.Run(configurator.GetMainConfig().GetListen()); nil != err {
		log.Fatalf("启动服务器失败: %v\n", err)
	}
}

const systemdConfigTemplate = `[Unit]
Description=V2ray Helper Service
Documentation=https://github.com/Luna-CY/v2ray-helper
After=network.target nss-lookup.target

[Service]
Type=simple
ExecStart=%v/v2ray-helper -home-dir %v
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target`

// installAsService 安装为系统服务
func installAsService(enable bool) {
	configFile, err := os.OpenFile("/etc/systemd/system/v2ray-helper.service", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		log.Fatalf("安装为系统服务失败: %v\n", err)
	}
	defer configFile.Close()

	if _, err := configFile.WriteString(fmt.Sprintf(systemdConfigTemplate, runtime.GetRootPath(), runtime.GetRootPath())); nil != err {
		log.Fatalf("安装为系统服务失败: %v\n", err)
	}

	log.Println("安装成功")

	if enable {
		_, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/v2ray-helper.service /etc/systemd/system/multi-user.target.wants/v2ray-helper.service").Output()
		if nil != err {
			log.Fatalf("设为开机自启失败: %v\n", err)
		}
	}

	log.Println("设为开机自启成功")
}

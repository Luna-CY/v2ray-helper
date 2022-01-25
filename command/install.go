package command

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/software/caddy"
	"github.com/Luna-CY/v2ray-helper/common/software/nginx"
	"github.com/Luna-CY/v2ray-helper/common/software/vhelper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "将服务安装到系统",
	Args:  cobra.NoArgs,
	Run:   install,
}

func init() {
	installCommand.Flags().BoolVar(&https, "https", false, "启用HTTPS协议，启用HTTPS需要申请HTTPS证书，指定此参数时必须提供 host 参数")
	installCommand.Flags().StringVar(&host, "host", "", "用于申请HTTPS证书的域名，设置 https 参数必须提供")
	installCommand.Flags().BoolVar(&enable, "enable", false, "设置为开机自启动")
}

var https, enable bool

const systemdConfigTemplate = `[Unit]
Description=V2ray Helper Service
Documentation=https://github.com/Luna-CY/v2ray-helper
After=network.target nss-lookup.target

[Service]
Type=simple
ExecStart=%v run --home %v
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target`

func install(*cobra.Command, []string) {
	if https && "" == host {
		log.Fatalln("启用HTTPS时必须提供用于申请证书的域名")
	}

	configFile, err := os.OpenFile("/etc/systemd/system/v2ray-helper.service", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if nil != err {
		log.Fatalf("安装为系统服务失败: %v\n", err)
	}
	defer configFile.Close()

	rootPath := viper.GetString(configurator.KeyRootPath)
	if _, err := configFile.WriteString(fmt.Sprintf(systemdConfigTemplate, filepath.Join(rootPath, "v2ray-helper"), rootPath)); nil != err {
		log.Fatalf("安装为系统服务失败: %v\n", err)
	}

	if enable {
		if _, err := exec.Command("sh", "-c", "ln -sf /etc/systemd/system/v2ray-helper.service /etc/systemd/system/multi-user.target.wants/v2ray-helper.service").Output(); nil != err {
			log.Fatalf("设为开机自启失败: %v\n", err)
		}
	}

	if https {
		viper.Set(configurator.KeyServerHttpsEnable, true)
		viper.Set(configurator.KeyServerHttpsHost, host)

		if err := viper.WriteConfig(); nil != err {
			log.Fatalln(err)
		}

		// 如果有Nginx服务器并且已启动，那么停止Nginx，否则Caddy无法启动
		nginxIsRunning, err := nginx.IsRunning()
		if nil != err {
			log.Fatalln(err)
		}

		if nginxIsRunning {
			if err := nginx.Stop(); nil != err {
				log.Fatalln(err)
			}
		}

		if err := nginx.Disable(); nil != err {
			log.Fatalln(err)
		}

		caddyIsRunning, err := caddy.IsRunning()
		if nil != err {
			log.Fatalln(err)
		}

		if caddyIsRunning {
			if err := caddy.Stop(); nil != err {
				log.Fatalln(err)
			}
		}

		vHelperIsRunning, err := vhelper.IsRunning()
		if nil != err {
			log.Fatalln(err)
		}

		if vHelperIsRunning {
			if err := vhelper.Stop(); nil != err {
				log.Fatalln(err)
			}
		}

		if _, err := certificate.GetManager().IssueNew(host, viper.GetString(configurator.KeyHttpsIssueEmail)); nil != err {
			log.Fatalln(err)
		}
	}

	if err := vhelper.Start(); nil != err {
		log.Fatalln(err)
	}

	log.Println("安装成功，服务已启动")
}

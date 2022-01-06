package command

import (
	"context"
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

func init() {
	cmd := &cobra.Command{
		Use:   "cert",
		Short: "证书管理工具",
		Args:  cobra.NoArgs,
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "申请新的证书",
		Long:  "申请新的证书，申请之前必须确保服务器的80与443端口没有被占用，否则将无法完成证书申请的挑战任务",
		Args:  cobra.NoArgs,
		Run:   issue,
	}

	issueCmd.Flags().StringVar(&home, "home", "", "运行主目录，默认为服务命令所在目录")
	issueCmd.Flags().StringVar(&host, "host", "", "需要申请证书的域名")

	renewCmd := &cobra.Command{
		Use:   "renew",
		Short: "续期证书",
		Long:  "续期已申请的证书，续期之前必须确保服务器的80与443端口没有被占用，否则将无法完成证书续期的挑战任务",
		Args:  cobra.NoArgs,
		Run:   renew,
	}

	renewCmd.Flags().StringVar(&home, "home", "", "运行主目录，默认为服务命令所在目录")
	renewCmd.Flags().StringVar(&host, "host", "", "需要续期的域名")

	cmd.AddCommand(issueCmd)
	cmd.AddCommand(renewCmd)

	command.AddCommand(cmd)
}

func issue(*cobra.Command, []string) {
	home = filepath.Clean(strings.TrimSpace(home))
	host = strings.TrimSpace(host)
	rootAbsPath := runtime.AbsRootPath(home)

	if "" == host {
		log.Fatalln("域名不能为空")
	}

	if err := configurator.Init(rootAbsPath); nil != err {
		log.Fatalf("无法初始化配置参数: %v\n", err)
	}

	// logger组件需要在其他组件之前初始化
	if err := logger.Init(rootAbsPath); nil != err {
		log.Fatalf("初始化日志失败: %v\n", err)
	}

	if err := certificate.Init(context.Background()); nil != err {
		log.Fatalf("初始化证书管理器失败: %v\n", err)
	}

	if certificate.GetManager().CheckExists(host) {
		log.Fatalln("该域名证书已存在，如域名过期请使用 renew 命令")
	}

	cert, err := certificate.GetManager().IssueNew(host, viper.GetString(configurator.KeyHttpsIssueEmail))
	if nil != err {
		log.Fatalf("申请域名证书失败: %v\n", err)
	}

	log.Println("申请域名证书成功")
	log.Printf("证书位置: %v\n", cert.GetCertificateFilePath())
	log.Printf("证书私钥位置: %v\n", cert.GetPrivateKeyFilePath())
}

func renew(*cobra.Command, []string) {
	home = filepath.Clean(strings.TrimSpace(home))
	host = strings.TrimSpace(host)
	rootAbsPath := runtime.AbsRootPath(home)

	if "" == host {
		log.Fatalln("域名不能为空")
	}

	if err := configurator.Init(rootAbsPath); nil != err {
		log.Fatalf("无法初始化配置参数: %v\n", err)
	}

	// logger组件需要在其他组件之前初始化
	if err := logger.Init(rootAbsPath); nil != err {
		log.Fatalf("初始化日志失败: %v\n", err)
	}

	if err := certificate.Init(context.Background()); nil != err {
		log.Fatalf("初始化证书管理器失败: %v\n", err)
	}

	if !certificate.GetManager().CheckExists(host) {
		log.Fatalln("该域名证书不存在，如需申请证书请使用 issue 命令")
	}

	if err := certificate.GetManager().Renew(host, viper.GetString(configurator.KeyHttpsIssueEmail)); nil != err {
		log.Fatalf("域名续期失败: %v\n", err)
	}

	log.Println("域名续期成功")
}

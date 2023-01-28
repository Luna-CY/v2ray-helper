package command

import (
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var certCommand = &cobra.Command{
	Use:   "cert",
	Short: "证书管理工具",
	Args:  cobra.NoArgs,
}

func init() {
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "申请新的证书",
		Long:  "申请新的证书，申请之前必须确保服务器的80与443端口没有被占用，否则将无法完成证书申请的挑战任务",
		Args:  cobra.NoArgs,
		Run:   issue,
	}

	issueCmd.Flags().StringVar(&host, "host", "", "需要申请证书的域名")

	renewCmd := &cobra.Command{
		Use:   "renew",
		Short: "续期证书",
		Long:  "续期已申请的证书，续期之前必须确保服务器的80与443端口没有被占用，否则将无法完成证书续期的挑战任务",
		Args:  cobra.NoArgs,
		Run:   renew,
	}

	renewCmd.Flags().StringVar(&host, "host", "", "需要续期的域名")

	certCommand.AddCommand(issueCmd)
	certCommand.AddCommand(renewCmd)
}

func issue(*cobra.Command, []string) {
	host = strings.TrimSpace(host)

	if "" == host {
		log.Fatalln("域名不能为空")
	}

	if certificate.GetManager().CheckExists(host) {
		log.Fatalln("该域名证书已存在，如域名过期请使用 renew 命令")
	}

	cert, err := certificate.GetManager().IssueNew(host)
	if nil != err {
		log.Fatalf("申请域名证书失败: %v\n", err)
	}

	log.Println("申请域名证书成功")
	log.Printf("证书位置: %v\n", cert.GetCertificateFilePath())
	log.Printf("证书私钥位置: %v\n", cert.GetPrivateKeyFilePath())
}

func renew(*cobra.Command, []string) {
	host = strings.TrimSpace(host)

	if "" == host {
		log.Fatalln("域名不能为空")
	}

	if !certificate.GetManager().CheckExists(host) {
		log.Fatalln("该域名证书不存在，如需申请证书请使用 issue 命令")
	}

	if err := certificate.GetManager().Renew(host); nil != err {
		log.Fatalf("域名续期失败: %v\n", err)
	}

	log.Println("域名续期成功")
}

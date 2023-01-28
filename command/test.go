package command

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/mail"
	"github.com/spf13/cobra"
	"log"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "测试工具集合",
	Args:  cobra.NoArgs,
}

func init() {
	ste := &cobra.Command{
		Use:   "send-test-email",
		Short: "发送测试邮件，测试邮件配置",
		Args:  cobra.NoArgs,
		Run:   sendTestEmail,
	}

	testCommand.AddCommand(ste)
}

func sendTestEmail(*cobra.Command, []string) {
	if err := mail.SendTestEmail(configurator.Configure.Mail.Notice.To); nil != err {
		log.Fatalf("发送测试邮件失败: %v\n", err)
	}

	log.Println("发送成功，请查收")
}

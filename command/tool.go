package command

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/mail"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	cmd := &cobra.Command{
		Use:   "tool",
		Short: "工具集合",
		Args:  cobra.NoArgs,
	}

	ste := &cobra.Command{
		Use:   "send-test-email",
		Short: "发送测试邮件，测试邮件配置",
		Args:  cobra.NoArgs,
		Run:   sendTestEmail,
	}

	cmd.AddCommand(ste)
	command.AddCommand(cmd)
}

func sendTestEmail(*cobra.Command, []string) {
	if err := mail.SendTestEmail(viper.GetString(configurator.KeyMailNoticeTo)); nil != err {
		log.Fatalf("发送测试邮件失败: %v\n", err)
	}

	log.Println("发送成功，请查收")
}

package mail

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"net/smtp"
	"strings"
)

// SendCertRenewFailEmail 发送证书续期失败的通知邮件
func SendCertRenewFailEmail(to, host string) error {
	if "" == to {
		return errors.New("接收地址不能为空")
	}

	address := strings.TrimSpace(viper.GetString(configurator.KeyMailSMTPServerAddress))
	port := viper.GetInt(configurator.KeyMailSMTPServerPort)
	user := strings.TrimSpace(viper.GetString(configurator.KeyMailSMTPUser))
	password := strings.TrimSpace(viper.GetString(configurator.KeyMailSMTPPassword))
	secret := strings.TrimSpace(viper.GetString(configurator.KeyMailSMTPSecret))

	if "" == address {
		return errors.New("邮件SMTP服务器配置错误，服务器地址不能为空")
	}

	mail := email.NewEmail()
	mail.From = user
	mail.To = []string{to}
	mail.Subject = "证书续期失败通知邮件"
	mail.HTML = []byte(fmt.Sprintf("<p>以下域名的证书续期失败，请及时处理</p><p>%v</p>", host))

	var auth smtp.Auth
	if "" != secret {
		auth = smtp.CRAMMD5Auth(user, secret)
	} else {
		auth = smtp.PlainAuth("", user, password, address)
	}

	if err := mail.Send(fmt.Sprintf("%v:%v", address, port), auth); nil != err {
		return err
	}

	return nil
}

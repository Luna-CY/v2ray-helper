package mail

import (
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/spf13/viper"
	"testing"
)

func TestSendCertRenewFailEmail(t *testing.T) {
	viper.Set(configurator.KeyMailSMTPServerAddress, "smtp.qq.com")
	viper.Set(configurator.KeyMailSMTPServerPort, 587)
	viper.Set(configurator.KeyMailSMTPUser, "your-email")
	viper.Set(configurator.KeyMailSMTPPassword, "replace-to-your-password")

	if err := SendCertRenewFailEmail("your-email", "test.luna.xin"); nil != err {
		t.Fatal(err)
	}
}

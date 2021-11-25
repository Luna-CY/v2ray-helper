package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/software/vhelper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

const (
	keyListen    = "listen"
	keyHttpsHost = "https-host"
	keyEmail     = "email"
)

type SaveMetaInfoForm struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SaveMetaInfo(c *gin.Context) {
	var body SaveMetaInfoForm
	if err := c.ShouldBind(&body); nil != err {
		response.Response(c, code.BadRequest, "无效请求", nil)

		return
	}

	body.Key = strings.TrimSpace(body.Key)
	body.Value = strings.TrimSpace(body.Value)

	restart := false

	switch body.Key {
	case keyListen:
		port, err := strconv.Atoi(body.Value)
		if nil != err {
			response.Response(c, code.BadRequest, "无法解析端口号", nil)

			return
		}

		if 80 == port || 443 == port || 1 > port || 65535 < port {
			response.Response(c, code.BadRequest, "无效端口号", nil)

			return
		}

		restart = true

		viper.Set(configurator.KeyServerPort, port)
	case keyHttpsHost:
		if "" != body.Value {
			if !certificate.GetManager().CheckExists(body.Value) {
				_, err := certificate.GetManager().IssueNew(body.Value, viper.GetString(configurator.KeyHttpsIssueEmail))
				if nil != err {
					logger.GetLogger().Errorln(err)
					response.Response(c, code.ServerError, "申请HTTPS证书失败，请稍后重试或联系管理员。详细错误请查看日志", nil)

					return
				}
			}

			viper.Set(configurator.KeyServerHttpsEnable, true)
			viper.Set(configurator.KeyServerHttpsHost, body.Value)
		} else {
			viper.Set(configurator.KeyServerHttpsEnable, false)
			viper.Set(configurator.KeyServerHttpsHost, "")
		}

		restart = true
	case keyEmail:
		_, err := mail.ParseAddress(body.Value)
		if nil != err {
			response.Response(c, code.BadRequest, "无效的邮箱地址", nil)

			return
		}

		viper.Set(configurator.KeyHttpsIssueEmail, body.Value)
	default:
		response.Response(c, code.BadRequest, "无效的数据请求", nil)

		return
	}

	if err := viper.WriteConfig(); nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.ServerError, "保存配置失败，请稍后重试或联系管理员。详细错误请查看日志", nil)

		return
	}

	response.Success(c, code.OK, nil)

	if restart && viper.GetBool(configurator.KeyServerRelease) {
		go func() {
			time.Sleep(3 * time.Second)

			_ = vhelper.ReStart()
		}()
	}

	return
}

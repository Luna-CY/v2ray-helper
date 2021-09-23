package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/certificate"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/Luna-CY/v2ray-helper/common/logger"
	"github.com/Luna-CY/v2ray-helper/common/runtime"
	"github.com/Luna-CY/v2ray-helper/common/software/vhelper"
	"github.com/gin-gonic/gin"
	"net/mail"
	"path/filepath"
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

		configurator.GetMainConfig().Listen = port
	case keyHttpsHost:
		if "" != body.Value {
			if !certificate.GetManager().CheckExists(body.Value) {
				_, err := certificate.GetManager().IssueNew(body.Value, configurator.GetMainConfig().Email)
				if nil != err {
					logger.GetLogger().Errorln(err)
					response.Response(c, code.ServerError, "申请HTTPS证书失败，请稍后重试或联系管理员。详细错误请查看日志", nil)

					return
				}
			}

			configurator.GetMainConfig().EnableHttps = true
			configurator.GetMainConfig().HttpsHost = body.Value
		} else {
			configurator.GetMainConfig().EnableHttps = false
			configurator.GetMainConfig().HttpsHost = ""
		}

		restart = true
	case keyEmail:
		_, err := mail.ParseAddress(body.Value)
		if nil != err {
			response.Response(c, code.BadRequest, "无效的邮箱地址", nil)

			return
		}

		configurator.GetMainConfig().Email = body.Value
	default:
		response.Response(c, code.BadRequest, "无效的数据请求", nil)

		return
	}

	if err := configurator.GetMainConfig().Save(filepath.Join(runtime.GetRootPath(), "config", configurator.GetMainConfig().GetFileName())); nil != err {
		logger.GetLogger().Errorln(err)
		response.Response(c, code.ServerError, "保存端口失败，请稍后重试或联系管理员。详细错误请查看日志", nil)

		return
	}

	response.Success(c, code.OK, nil)

	if restart && configurator.GetMainConfig().GinReleaseMode {
		go func() {
			time.Sleep(3 * time.Second)

			_ = vhelper.ReStart()
		}()
	}

	return
}

package middleware

import (
	"errors"
	"github.com/Luna-CY/v2ray-helper/code"
	"github.com/Luna-CY/v2ray-helper/configurator"
	"github.com/Luna-CY/v2ray-helper/logger"
	"github.com/Luna-CY/v2ray-helper/response"
	"github.com/Luna-CY/v2ray-helper/util"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

const secret = "v2ray-jwt-secret"
const JwtIdentityKey = "userinfo"

// AuthForm 验证表单
type AuthForm struct {
	Key string `json:"key" binding:"required"`
}

func GetJWT() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       secret,
		Key:         []byte(secret),
		MaxRefresh:  30 * time.Minute,
		IdentityKey: JwtIdentityKey,
		IdentityHandler: func(context *gin.Context) interface{} {
			cl := jwt.ExtractClaims(context)

			return cl[JwtIdentityKey].(string)
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{JwtIdentityKey: data.(string)}
		},
		Authenticator: func(context *gin.Context) (interface{}, error) {
			var body AuthForm

			if err := context.ShouldBind(&body); nil != err {
				logger.GetLogger().Errorf("绑定数据失败: %v\n", err)

				return nil, errors.New("无效的数据请求")
			}

			if util.Md5(configurator.GetMainConfig().Key) != strings.TrimSpace(body.Key) {
				return nil, errors.New("无效口令")
			}

			context.Set(JwtIdentityKey, "success")

			return "success", nil
		},
		LoginResponse: func(context *gin.Context, httpCode int, token string, t time.Time) {
			_, ok := context.Get(JwtIdentityKey)
			if !ok {
				logger.GetLogger().Errorln("无法获取登录的用户信息")
				response.Response(context, code.ServerError, "服务器内部错误", nil)

				return
			}

			response.Success(context, code.OK, &gin.H{
				"token":   token,
				"expired": t.Unix(),
			})
		},
		Unauthorized: func(context *gin.Context, httpCode int, message string) {
			if message == jwt.ErrExpiredToken.Error() {
				response.Response(context, code.TokenIsExpired, message, nil)

				return
			}

			response.Response(context, code.UnAuthorized, message, nil)
		},
		LogoutResponse: func(context *gin.Context, httpCode int) {
			response.Success(context, code.OK, nil)
		},
		RefreshResponse: func(context *gin.Context, hc int, token string, t time.Time) {
			response.Response(context, code.BadRequest, "不管你是谁，请停止你的行为", nil)
		},
	})
}

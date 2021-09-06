package middleware

import (
	"encoding/json"
	"errors"
	"gitee.com/Luna-CY/v2ray-subscription/code"
	"gitee.com/Luna-CY/v2ray-subscription/database/model"
	"gitee.com/Luna-CY/v2ray-subscription/dataservice"
	"gitee.com/Luna-CY/v2ray-subscription/logger"
	"gitee.com/Luna-CY/v2ray-subscription/response"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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

			key := new(model.Key)
			if err := json.Unmarshal([]byte(cl[JwtIdentityKey].(string)), key); nil != err {
				return nil
			}

			return key
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			d, _ := json.Marshal(data)

			return jwt.MapClaims{JwtIdentityKey: string(d)}
		},
		Authenticator: func(context *gin.Context) (interface{}, error) {
			var authForm AuthForm

			if err := context.ShouldBind(&authForm); nil != err {
				logger.GetLogger().Errorf("绑定数据失败: %v\n", err)

				return nil, errors.New("无效的数据请求")
			}

			key := new(model.Key)
			if err := dataservice.GetBaseService().TakeByCondition(key, nil, "key = ?", authForm.Key); nil != err {
				if gorm.ErrRecordNotFound == err {
					return nil, errors.New("无效口令")
				}

				return nil, errors.New("服务器内部错误，请稍后再试")
			}

			context.Set(JwtIdentityKey, key)

			return key, nil
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
			if http.StatusOK == hc {
				response.Success(context, code.OK, &gin.H{
					"token":   token,
					"expired": t.Unix(),
				})

				return
			}

			response.Response(context, code.ServerError, "服务器内部错误", nil)
		},
	})
}

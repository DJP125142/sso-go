package router

import (
	"github.com/gin-gonic/gin"
	"sso-go/controller"
	"sso-go/middlewares"
)

func AccountRouter(Router *gin.RouterGroup) {
	AccountRouter := Router.Group("account")
	{
		// 发送邮箱验证码
		AccountRouter.POST("send_emial_code", controller.SendValidateCode)
		// 注册
		AccountRouter.POST("register", controller.Register)
		// 登录
		AccountRouter.POST("login", controller.Login)
		// 获取用户信息
		AccountRouter.GET("user", middlewares.JWTAuth(), controller.UserInfo)
		// 创建授权code
		AccountRouter.POST("create_code", controller.CreateCode)
		// 外部客户端拿code换取token
		AccountRouter.POST("get_token_by_code", controller.GetTokenByCode)
	}
}

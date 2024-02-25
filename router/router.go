package router

import (
	"github.com/gin-gonic/gin"
	"sso-go/controller"
	"sso-go/middlewares"
)

func AccountRouter(Router *gin.RouterGroup) {
	AccountRouter := Router.Group("account")
	{
		AccountRouter.POST("send_emial_code", controller.SendValidateCode)
		AccountRouter.POST("register", controller.Register)
		AccountRouter.POST("login", controller.Login)
		AccountRouter.GET("user", middlewares.JWTAuth(), controller.UserInfo)
	}
}

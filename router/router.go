package router

import (
	"github.com/gin-gonic/gin"
	"sso-go/controller"
)

func AccountRouter(Router *gin.RouterGroup) {
	AccountRouter := Router.Group("account")
	{
		AccountRouter.GET("test", controller.Test)
	}
}

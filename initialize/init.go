package initialize

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sso-go/config"
	"sso-go/global"
	"sso-go/router"
)

/*
* 初始化配置项
 */
func InitConfig() {
	// 实例化viper
	v := viper.New()
	v.SetConfigName("env")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 声明一个ServerConfig类型的实例
	serverConfig := config.ServerConfig{}
	// 给serverConfig初始值
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	// 传递给全局变量
	global.Settings = serverConfig
	color.Blue("initConfig", global.Settings.LogsAddress)
}

/*
* 初始化路由
 */
func InitRouters() *gin.Engine {
	Router := gin.Default()
	// 路由分组
	ApiGroup := Router.Group("/v1/")
	router.AccountRouter(ApiGroup) // 注册AccountRouter组路由
	return Router
}

package main

import (
	"fmt"
	"sso-go/global"
	"sso-go/initialize"
)

func main() {
	// 1.初始化yaml配置
	initialize.InitConfig()
	// 2.初始化routers
	Router := initialize.InitRouters()
	// 3.初始化日志信息
	initialize.InitLogger()
	// 4.初始化语言翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	// 5.初始化mysql
	initialize.InitMysqlDB()
	// 6.初始化redis
	initialize.InitRedis()

	Router.Run(fmt.Sprintf(":%d", global.Settings.Port))
}

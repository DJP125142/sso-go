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

	Router.Run(fmt.Sprintf(":%d", global.Settings.Port))
}

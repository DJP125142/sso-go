package initialize

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sso-go/config"
	"sso-go/global"
	"sso-go/middlewares"
	"sso-go/router"
	"sso-go/utils"
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
	emailConfig := config.EmailConfig{}
	// 给serverConfig初始值
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	// 传递给全局变量
	global.Settings = serverConfig
	global.Email = emailConfig
	color.Blue("initConfig", global.Settings.LogsAddress)
}

/*
* 初始化路由
 */
func InitRouters() *gin.Engine {
	Router := gin.Default()
	// 加载自定义中间件
	Router.Use(middlewares.GinLogger(), middlewares.GinRecovery(true))
	// 路由分组
	ApiGroup := Router.Group("/v1/")
	router.AccountRouter(ApiGroup) // 注册AccountRouter组路由
	return Router
}

// InitLogger 初始化Logger
func InitLogger() {
	// 实例化zap配置
	cfg := zap.NewDevelopmentConfig()
	// 配置日志的输出地址
	cfg.OutputPaths = []string{
		fmt.Sprintf("%slog_%s.log", global.Settings.LogsAddress, utils.GetNowFormatTodayTime()),
		"stdout", // "stdout" 表示同时将日志输出到标准输出流（控制台）。这样就可以将日志同时输出到文件和控制台
	}
	// 创建logger实例
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	global.Lg = logger         // 注册到全局变量中
}

// 初始化Mysql
func InitMysqlDB() {
	mysqlInfo := global.Settings.MysqlInfo
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.Username, mysqlInfo.Password, mysqlInfo.Host,
		mysqlInfo.Port, mysqlInfo.Database)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	global.DB = db
}

// 初始化redis
func InitRedis() {
	addr := fmt.Sprintf("%s:%d", global.Settings.RedisInfo.Host, global.Settings.RedisInfo.Port)
	password := fmt.Sprintf("%s", global.Settings.RedisInfo.Password)
	//生成redis客户端
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	// 连接redis
	_, err := global.Redis.Ping().Result()
	if err != nil {
		color.Red("[InitRedis] 链接redis异常:")
		color.Yellow(err.Error())
	}
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"sso-go/global"
	"strings"
	"time"
)

// GinLogger 是一个gin中间件函数，用于记录请求日志信息
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求处理的开始时间
		start := time.Now()
		// 获取请求路径和请求参数
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 让请求继续往下走到后续的中间件或者controller业务逻辑
		c.Next()
		// 获取请求处理的耗时
		cost := time.Since(start)
		// 若response的状态码不是200为异常，记录异常信息
		if c.Writer.Status() != 200 {
			zap.L().Info(path,
				zap.Int("status", c.Writer.Status()),                                 // 记录状态码
				zap.String("method", c.Request.Method),                               // 记录请求方法
				zap.String("path", path),                                             // 记录请求路径
				zap.String("query", query),                                           // 记录请求参数
				zap.String("ip", c.ClientIP()),                                       // 记录请求IP
				zap.String("user-agent", c.Request.UserAgent()),                      // 记录请求用户代理
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 记录异常信息
				zap.Duration("cost", cost),                                           // 记录请求处理耗时
			)
		}
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	// 返回一个 gin.HandlerFunc 类型的函数作为中间件
	return func(c *gin.Context) {
		// 使用 defer 机制，当捕获到 panic 时执行相应的处理函数
		// defer是Go语言的一种语法机制，用于注册延迟调用。所谓延迟调用，就是在函数结束时按照defer语句的逆序进行执行，也就是先注册的defer语句最后执行，后注册的defer语句最先执行。
		defer func() {
			// 如果存在 panic 错误，执行以下处理
			if err := recover(); err != nil {
				// 判断错误是否是连接异常引发的
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 记录 http 请求的信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 如果是连接异常，记录错误日志，但不中断请求处理
				if brokenPipe {
					global.Lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				// 如果需要记录堆栈信息，记录错误日志和堆栈信息
				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					// 否则只记录错误日志
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 中断请求的处理，返回 500 状态码
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 让请求继续往下走到后续的中间件或者 controller 业务逻辑
		c.Next()
	}
}

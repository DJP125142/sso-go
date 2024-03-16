package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso-go/dao"
	"sso-go/forms"
	"sso-go/global"
	"sso-go/middlewares"
	"sso-go/model"
	"sso-go/response"
	"sso-go/utils"
	"time"
)

// 注册接口
func Register(c *gin.Context) {
	// 初始化 RegisterForm 结构体
	registerParams := forms.RegisterForm{}
	// 使用 c.ShouldBind 函数将请求中的参数绑定到 RegisterForm 结构体上，如果出现错误，则将错误返回给客户端
	if err := c.ShouldBind(&registerParams); err != nil {
		// 统一处理异常
		utils.HandleValidatorError(c, err)
		return
	}

	// 验证邮箱验证码
	emailCodeKey := fmt.Sprintf("EmailCode:%s", registerParams.Email)
	if registerParams.Code != global.Redis.Get(emailCodeKey).Val() {
		response.Err(c, http.StatusOK, 400, "邮箱验证码错误", nil)
		return
	}

	// 验证邮箱或昵称是否注册
	hasEmail := dao.HasUser(registerParams.Email)
	hasName := dao.HasUser(registerParams.Username)
	if hasEmail {
		response.Err(c, http.StatusOK, 400, "该邮箱已注册", nil)
		return
	}
	if hasName {
		response.Err(c, http.StatusOK, 400, "该昵称已注册", nil)
		return
	}

	// 生成加密密码
	hashPwd := utils.HashAndSalt(registerParams.PassWord)

	// 创建用户
	user := model.User{
		Name:            registerParams.Username,
		Email:           registerParams.Email,
		EmailVerifiedAt: utils.GetNowFormatTime(),
		Password:        hashPwd,
	}
	result := global.DB.Create(&user)
	if result.Error != nil {
		response.Err(c, 200, 500, "创建失败", result.Error.Error())
		return
	}

	data := map[string]interface{}{
		"user_id": user.ID,
	}
	// 如果请求参数绑定成功，则返回状态码为 200 和成功信息给客户端
	response.Success(c, 200, "success", data)
	return
}

// 登录接口
func Login(c *gin.Context) {
	// 初始化 RegisterForm 结构体
	loginParams := forms.LoginForm{}
	// 使用 c.ShouldBind 函数将请求中的参数绑定到 RegisterForm 结构体上，如果出现错误，则将错误返回给客户端
	if err := c.ShouldBind(&loginParams); err != nil {
		// 统一处理异常
		utils.HandleValidatorError(c, err)
		return
	}

	// 查询是否有该用户
	user, ok, msg := dao.GetUserInfoByPw(loginParams.Username, loginParams.PassWord)
	if !ok {
		response.Err(c, 401, 401, msg, "")
		return
	}

	// 处理下默认头像
	if user.HeadUrl == "" {
		user.HeadUrl = "http://account.djp.org.cn/static/img/head_default.1e4f19e.png"
	}

	// 登录成功创建token
	token := utils.CreateToken(c, user.ID, user.Name, user.Email, user.HeadUrl)
	userInfoMap := HandleUserModelToMap(user)
	userInfoMap["token"] = token

	response.Success(c, 200, "success", userInfoMap)
}

// 用户信息
func UserInfo(c *gin.Context) {
	c.Get("claims")
	claims, exists := c.Get("claims")
	if !exists {
		// 如果不存在，说明中间件没有设置claims
		response.Err(c, 401, 401, "未登录", "")
		return
	}

	// 进行类型断言，确保 claims 的类型是 jwt.Claims
	jwtClaims, ok := claims.(middlewares.CustomClaims)
	if !ok {
		response.Err(c, 401, 401, "未登录", "")
		return
	}

	userInfo := map[string]interface{}{
		"userId":   jwtClaims.ID,
		"username": jwtClaims.NickName,
		"email":    jwtClaims.Email,
	}
	response.Success(c, 200, "success", map[string]interface{}{
		"userInfo": userInfo,
	})
}

// 发送邮箱验证码
func SendValidateCode(c *gin.Context) {
	email := c.Query("email")
	if utils.IsEmail(email) {
		response.Err(c, 401, 401, "email格式错误", "")
		return
	}
	emails := []string{email}
	// 发送邮件
	vCode, err := utils.SendEmailValidate(emails)
	if err != nil {
		response.Err(c, 200, 500, "验证码发送失败", err.Error())
		return
	}
	// 验证码存入redis，有效期5分钟
	emailCodeKey := fmt.Sprintf("EmailCode:%s", email)
	global.Redis.Set(emailCodeKey, vCode, time.Minute*5)

	response.Success(c, 200, "success", nil)
	return
}

func HandleUserModelToMap(user *model.User) map[string]interface{} {
	userItemMap := map[string]interface{}{
		"id":       user.ID,
		"username": user.Name,
		"head_url": user.HeadUrl,
		"email":    user.Email,
	}
	return userItemMap
}

// 获取临时授权码
func CreateCode(c *gin.Context) {
	cookie, _ := c.Request.Cookie("token")
	token := cookie.Value
	code := utils.GenerateCode()

	// code存入redis，有效期1分钟
	global.Redis.Set(code, token, time.Minute)
	response.Success(c, 200, "success", code)
}

// 根据code来换取token
func GetTokenByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Err(c, 401, 401, "code不得为空", "")
		return
	}
	token := global.Redis.Get(code).Val()

	// 校验token
	j := middlewares.NewJWT()
	claims, err := j.ParseToken(token)
	global.Lg.Info("claims info", zap.Any("claims", claims))
	if err != nil {
		response.Err(c, 401, 401, "fail", "")
		return
	}

	// 验证后清除redis
	global.Redis.Del(code)

	userInfo := map[string]interface{}{
		"userId":        claims.ID,
		"username":      claims.NickName,
		"head_url":      claims.HeadUrl,
		"token":         token,
		"expirein_time": claims.StandardClaims.ExpiresAt,
	}

	response.Success(c, 200, "success", userInfo)

}

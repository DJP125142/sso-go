package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"regexp"
	"sso-go/global"
	"sso-go/middlewares"
	"sso-go/response"
	"strings"
	"time"
)

// 获取一个当前日期格式的字符串
func GetNowFormatTodayTime() string {
	now := time.Now()
	format := "2006-01-02"
	formattedTime := now.Format(format)
	return formattedTime
}

// 获取当前时间戳的格式
func GetNowFormatTime() string {
	now := time.Now()
	format := "2006-01-02 15:04:05"
	formattedTime := now.Format(format)
	return formattedTime
}

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	// 利用反射拿到结构体tag含有mobile的key字段
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}

// 校验邮箱
func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	if result {
		return true
	} else {
		return false
	}
}

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(c *gin.Context, err error) {
	//如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		response.Err(c, http.StatusOK, 500, "参数校验错误", err.Error())
	}
	msg := removeTopStruct(errs.Translate(global.Trans))
	response.Err(c, http.StatusOK, 400, "参数校验错误", msg)
	return
}

// removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		// 从文本的逗号开始切分   处理后"mobile": "mobile为必填字段"  处理前: "PasswordLoginForm.mobile": "mobile为必填字段"
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 加密密码
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// 验证密码
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// 发送邮箱验证码，支持群发
func SendEmailValidate(em []string) (string, error) {
	e := email.NewEmail()
	//服务器相关的配置
	emailConfig := global.Settings.EmailInfo
	e.Subject = "邮箱验证"
	e.From = fmt.Sprintf("%s <%s>", emailConfig.SendName, emailConfig.SendEmail)
	e.To = em
	// 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	//设置文件发送的内容
	content := fmt.Sprintf(`
	<div>
		<div>
			尊敬的%s，您好！
		</div>
		<div style="padding: 8px 40px 8px 50px;">
			<p>您于 %s 提交的邮箱验证，本次验证码为<u><strong>%s</strong></u>，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
		</div>
		<div>
			<p>此邮箱为系统邮箱，请勿回复。</p>
		</div>
	</div>
	`, em[0], t, vCode)
	e.HTML = []byte(content)
	global.Lg.Info("log email info ", zap.Any("emailInfo", emailConfig))
	err := e.Send(emailConfig.Address, smtp.PlainAuth("", emailConfig.SendEmail, emailConfig.Password, emailConfig.Host))
	return vCode, err
}

// 生成jwt的token
func CreateToken(c *gin.Context, Id uint, NickName string, Email string, HeadUrl string) string {
	//生成token信息
	j := middlewares.NewJWT()
	claims := middlewares.CustomClaims{
		ID:       Id,
		NickName: NickName,
		Email:    Email,
		HeadUrl:  HeadUrl,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			// TODO 设置token过期时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, //token -->7天过期
			Issuer:    "test",
		},
	}
	//生成token
	token, err := j.CreateToken(claims)
	if err != nil {
		response.Success(c, 401, "token生成失败,重新再试", "test")
		return ""
	}
	return token
}

// 随机生成一个code码
func GenerateCode() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(token)
}

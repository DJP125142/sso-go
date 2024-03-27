package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
	"sso-go/global"
	"sso-go/response"
	"strings"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	ID       uint
	NickName string
	Email    string
	HeadUrl  string
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 头部
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			response.Err(c, http.StatusOK, 401, "请登录", "")
			global.Lg.Info("jwt鉴权失败401：", zap.Any("error:", "没有Authorization"))
			c.Abort()
			return
		}
		// parseToken 解析token包含的信息
		token := ExtractTokenFromHeader(authorization)
		fmt.Println(token)
		j := NewJWT()
		claims, err := j.ParseToken(token)
		jwt.TimeFunc = time.Now
		if err != nil {
			if err == TokenExpired {
				if err == TokenExpired {
					response.Err(c, http.StatusOK, 401, "授权已过期", "")
					c.Abort()
					return
				}
			}
			response.Err(c, http.StatusOK, 401, "未登陆", "")
			c.Abort()
			return
		}
		// gin的上下文记录claims和userId的值
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// 辅助函数：从 Authorization 头部中提取 Token
func ExtractTokenFromHeader(authHeader string) string {
	// Token 应该以 "Bearer " 前缀开始，因此我们可以简单地删除前缀以获取 Token
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.Settings.JWTKey.SigningKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

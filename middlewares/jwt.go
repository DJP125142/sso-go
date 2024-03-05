package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"sso-go/global"
	"sso-go/response"
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
		// 我们这里jwt鉴权取头部信息 token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			response.Err(c, http.StatusUnauthorized, 401, "请登录", err.Error())
			c.Abort()
			return
		}
		if cookie.Value == "" {
			response.Err(c, http.StatusUnauthorized, 401, "请登录", "")
			c.Abort()
			return
		}
		token := cookie.Value
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		jwt.TimeFunc = time.Now
		if err != nil {
			if err == TokenExpired {
				if err == TokenExpired {
					response.Err(c, http.StatusUnauthorized, 401, "授权已过期", "")
					c.Abort()
					return
				}
			}
			response.Err(c, http.StatusUnauthorized, 401, "未登陆", "")
			c.Abort()
			return
		}
		fmt.Println(c)
		// gin的上下文记录claims和userId的值
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
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

package controller

import (
	"github.com/gin-gonic/gin"
	"sso-go/response"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	string   `json:"password"`
}

func Test(c *gin.Context) {
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	response.Success(c, 200, "success", map[string]interface{}{
		"data": "pong",
	})
}

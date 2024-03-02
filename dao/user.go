package dao

import (
	"go.uber.org/zap"
	"sso-go/global"
	"sso-go/model"
	"sso-go/utils"
)

var user model.User

// 用户是否存在
func HasUser(nameOrEmail string) bool {
	var whereMap map[string]interface{}
	if utils.IsEmail(nameOrEmail) {
		whereMap = map[string]interface{}{"email": nameOrEmail}
	} else {
		whereMap = map[string]interface{}{"name": nameOrEmail}
	}
	rows := global.DB.Limit(1).Where(whereMap).First(&user)
	if rows.RowsAffected < 1 {
		return false
	}
	return true
}

// UsernameFindUserInfo 通过username找到用户信息
func GetUserInfoByPw(username string, password string) (*model.User, bool, string) {
	var whereMap map[string]interface{}
	if utils.IsEmail(username) {
		whereMap = map[string]interface{}{"email": username}
	} else {
		whereMap = map[string]interface{}{"name": username}
	}
	global.Lg.Info("login log", zap.Any("whereMap", whereMap))
	rows := global.DB.Limit(1).Where(whereMap).First(&user)
	if rows.RowsAffected < 1 {
		return &user, false, "该用户未注册"
	}
	// 校验密码
	verifyPassword := utils.ComparePasswords(user.Password, password)
	if !verifyPassword {
		return &user, false, "密码验证失败"
	}
	return &user, true, "登录成功"
}

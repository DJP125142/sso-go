package forms

type RegisterForm struct {
	// 用户名
	Username string `form:"name" json:"name" binding:"required,min=2,max=20"`
	// 邮箱
	Email string `form:"email" json:"email" binding:"required,email"`
	// 邮箱验证码
	Code string `form:"code" json:"code" binding:"required,len=6"`
	// 密码
	PassWord string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

type LoginForm struct {
	// 用户名
	Username string `form:"name" json:"name" binding:"required,min=2,max=20"`
	// 密码
	PassWord string `form:"password" json:"password" binding:"required,min=6,max=20"`
}

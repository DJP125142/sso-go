package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sso-go/config"
)

var (
	Settings config.ServerConfig
	Lg       *zap.Logger
	Trans    ut.Translator
	Email    config.EmailConfig
	DB       *gorm.DB
	Redis    *redis.Client
)

package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"webfw/configs"
	"webfw/pkg/redis"
)

var (
	L *zap.Logger
	DB *gorm.DB
	Redis *redis.Redis
	C *configs.Config
)

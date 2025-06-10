package global

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"server/config"
)

var (
	DB     *gorm.DB
	REDIS  *redis.Client
	CONFIG config.Server
	VP     *viper.Viper
	CASBIN *casbin.Enforcer
	TIMER  sync.Map
	lock   sync.RWMutex
)

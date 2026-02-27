package redis

import (
	"fmt"
	"sync"

	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

var (
	once          sync.Once
	redisInstance *Redis
)

func NewRedisClient(conf *config.Config) *Redis {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port),
			DB:       0,
			Protocol: 2,
		})

		redisInstance = &Redis{
			Client: rdb,
		}
	})
	return redisInstance
}

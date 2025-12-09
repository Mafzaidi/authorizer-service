package redis

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"localdev.me/authorizer/config"
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
			Password: conf.Redis.Password,
			DB:       0,
			Protocol: 2,
		})

		redisInstance = &Redis{
			Client: rdb,
		}
	})
	return redisInstance
}

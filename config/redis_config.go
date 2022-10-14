package config

import (
	"sync"

	"github.com/go-redis/redis"
)

type RedisConfigMode int

const (
	RedisConfigModeOptions RedisConfigMode = 0
	RedisConfigModeRaw     RedisConfigMode = 1
)

type RedisConfig struct {
	Mode RedisConfigMode

	Options *redis.Options
	Client  *redis.Client

	once sync.Once
}

func (c *RedisConfig) InitClient() redis.Cmdable {
	c.once.Do(func() {
		if c.Mode == RedisConfigModeOptions {
			c.Client = redis.NewClient(c.Options)
		}
	})
	return c.Client
}

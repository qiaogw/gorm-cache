package config

import (
	"sync"

	"gopkg.in/redis.v4"
)

type RedisConfigMode int

const (
	RedisConfigModeOptions RedisConfigMode = 0
	RedisConfigModeRaw     RedisConfigMode = 1
)

type RedisConfig struct {
	Mode RedisConfigMode

	Options *redis.Options
	Client  redis.Cmdable

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

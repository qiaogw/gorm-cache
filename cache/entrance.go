package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/qiaogw/gorm-cache/config"
)

func NewGorm2Cache(cacheConfig *config.CacheConfig) (*Gorm2Cache, error) {
	if cacheConfig == nil {
		return nil, fmt.Errorf("you pass a nil config")
	}
	cache := &Gorm2Cache{
		Config: cacheConfig,
	}
	err := cache.Init()
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func NewRedisConfigWithOptions(options *redis.Options) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:    config.RedisConfigModeOptions,
		Options: options,
	}
}

// NewRedisConfigWithClient 新
func NewRedisConfigWithClient(client *redis.Client) *config.RedisConfig {
	return &config.RedisConfig{
		Mode:   config.RedisConfigModeRaw,
		Client: client,
	}
}

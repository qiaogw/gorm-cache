package data_layer

import (
	"context"

	"github.com/Pacific73/gorm-cache/config"
	"github.com/Pacific73/gorm-cache/util"
	"github.com/go-redis/redis"
)

type RedisLayer struct {
	client    *redis.Client
	ttl       int64
	logger    config.LoggerInterface
	keyPrefix string

	batchExistSha string
	cleanCacheSha string
}

func (r *RedisLayer) Init(conf *config.CacheConfig, prefix string) error {
	if conf.RedisConfig.Mode == config.RedisConfigModeOptions {
		r.client = redis.NewClient(conf.RedisConfig.Options)
	} else {
		r.client = conf.RedisConfig.Client
	}

	r.ttl = conf.CacheTTL
	r.logger = conf.DebugLogger
	r.keyPrefix = prefix
	return r.initScripts()
}

func (r *RedisLayer) initScripts() error {
	batchKeyExistScript := `
		for idx, val in pairs(KEYS) do
			local exists = redis.call('EXISTS', val)
			if exists == false then
				return false
			end
		end
		return true`

	cleanCacheScript := `
		local keys = redis.call('keys', ARGV[1])
		for i=1,#keys,5000 do 
			redis.call('del', 'defaultKey', unpack(keys, i, math.min(i+4999, #keys)))
		end
		return 1`

	result := r.client.ScriptLoad(batchKeyExistScript)
	if result.Err() != nil {
		r.logger.CtxError(context.Background(), "[initScripts] init script 1 error: %v", result.Err())
		return result.Err()
	}
	r.batchExistSha = result.String()

	result = r.client.ScriptLoad(cleanCacheScript)
	if result.Err() != nil {
		r.logger.CtxError(context.Background(), "[initScripts] init script 2 error: %v", result.Err())
		return result.Err()
	}
	r.cleanCacheSha = result.String()
}

func (r *RedisLayer) CleanCache(ctx context.Context) error {
	result := r.client.EvalSha(r.cleanCacheSha, []string{"0"}, r.keyPrefix+":*")
	if result.Err() != nil {
		r.logger.CtxError(ctx, "[CleanCache] clean cache error: %v", result.Err())
		return result.Err()
	}
	return nil
}

func (r *RedisLayer) BatchKeyExist(ctx context.Context, keys []string) (bool, error) {
	result := r.client.EvalSha(r.batchExistSha, keys)
	if result.Err() != nil {
		r.logger.CtxError(ctx, "[BatchKeyExist] eval script error: %v", result.Err())
		return false, result.Err()
	}
	return result.Bool()
}

func (r *RedisLayer) KeyExists(ctx context.Context, key string) (bool, error) {
	result := r.client.Exists(key)
	if result.Err() != nil {
		r.logger.CtxError(ctx, "[KeyExists] exists error: %v", result.Err())
		return false, result.Err()
	}
	if result.Val() == 1 {
		return true, nil
	}
	return false, nil
}

func (r *RedisLayer) GetValue(ctx context.Context, key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisLayer) BatchGetValues(ctx context.Context, keys []string) ([]string, error) {
	result := r.client.MGet(keys...)
	if result.Err() != nil {
		r.logger.CtxError(ctx, "[BatchGetValues] mget error: %v", result.Err())
		return nil, result.Err()
	}
	slice := result.Val()
	strs := make([]string, 0, len(slice))
	for _, obj := range slice {
		strs = append(strs, obj.(string))
	}
	return strs, nil
}

func (r *RedisLayer) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {

}

func (r *RedisLayer) DeleteKey(ctx context.Context, key string) error {

}

func (r *RedisLayer) BatchDeleteKeys(ctx context.Context, keys []string) error {

}

func (r *RedisLayer) BatchSetKeys(ctx context.Context, kvs []util.Kv) error {

}

func (r *RedisLayer) SetKey(ctx context.Context, kv util.Kv) error {

}

package persistence

import (
	"context"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisPersistence struct {
	logger      adapters2.LoggerAdapter
	redisClient adapters.RedisAdapter
	ctx         context.Context
	prefix      string
	keyTTL      time.Duration
}

// RedisPersistence constructor for class.
func RedisPersistence(logger adapters2.LoggerAdapter, redisClient adapters.RedisAdapter) *redisPersistence {
	return &redisPersistence{
		logger:      logger,
		redisClient: redisClient,
		ctx:         context.Background(),
		prefix:      properties.Properties().Cache.KeyPrefix,
		keyTTL:      properties.Properties().Cache.KeyTTL,
	}
}

// Lock will set the key on cache with TTL active. Returns error if a problem occurs while trying to persist on cache.
func (r *redisPersistence) Lock(key string) custom_error.BaseErrorAdapter {
	r.logger.Info("Lock started", key)

	redisClient, err := r.redisClient.Open()
	if err != nil {
		return r.abort(err, "Error while trying to open redis connection", false, false)
	}

	_, err = redisClient.Get(r.ctx, r.prefix+key).Result()
	if err == nil {
		return r.abort(err, "Key is already locked", false, true)
	} else if err != nil && err != redis.Nil {
		return r.abort(err, "Error while trying to get redis key", false, true)
	}

	_, err = redisClient.Set(r.ctx, r.prefix+key, key, r.keyTTL).Result()
	if err != nil {
		return r.abort(err, "Error while trying to set redis key", false, true)
	}

	err = r.redisClient.Close()
	if err != nil {
		r.logger.Warning(err, "Could not close Redis connection")
	}

	r.logger.Info("Lock finished", key)
	return nil
}

// Unlock will remove the key from cache. Returns error if a problem occurs while trying to delete from cache.
func (r *redisPersistence) Unlock(key string) custom_error.BaseErrorAdapter {
	r.logger.Info("Unlock started", key)

	redisClient, err := r.redisClient.Open()
	if err != nil {
		return r.abort(err, "Error while trying to open redis connection", true, false)
	}

	_, err = redisClient.Del(r.ctx, r.prefix+key).Result()
	if err != nil {
		return r.abort(err, "Error while trying to delete redis key", true, true)
	}

	err = r.redisClient.Close()
	if err != nil {
		r.logger.Warning(err, "Could not close Redis connection")
	}

	r.logger.Info("Unlock finished", key)
	return nil
}

func (r *redisPersistence) abort(err error, message string, locked bool, closeConn bool) custom_error.BaseErrorAdapter {
	if closeConn {
		closeErr := r.redisClient.Close()
		if closeErr != nil {
			r.logger.Warning(err, "Could not close Redis connection")
		}
	}

	redisPersistenceLockError := exceptions.RedisPersistenceLockError(err, message, locked)
	r.logger.Error(redisPersistenceLockError, "Unlock failed: "+message)
	return redisPersistenceLockError
}

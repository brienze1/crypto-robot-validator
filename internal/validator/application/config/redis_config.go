package config

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisClient struct {
	secretsManager adapters.SecretsManagerServiceAdapter
	client         *redis.Client
	cacheConfig    *dto.RedisSecrets
}

func RedisClient(secretsManager adapters.SecretsManagerServiceAdapter) *redisClient {
	return &redisClient{
		secretsManager: secretsManager,
	}
}

func (r *redisClient) Open() (*redis.Client, error) {
	if r.cacheConfig == nil {
		cacheConfig, err := r.secretsManager.GetSecret(properties.Properties().Aws.SecretsManager.CacheSecretName)
		if err != nil {
			panic(err)
		}

		r.cacheConfig = cacheConfig
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:         r.cacheConfig.Address,
		Username:     r.cacheConfig.User,
		Password:     r.cacheConfig.Password,
		DB:           r.cacheConfig.DatabaseId,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	_, err := r.client.Ping(r.client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return r.client, nil
}

func (r *redisClient) Close() error {
	return r.client.Close()
}

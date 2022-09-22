package persistence

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

type redisPersistence struct {
}

// RedisPersistence constructor for class.
func RedisPersistence() *redisPersistence {
	return &redisPersistence{}
}

// Lock will set the key on cache with TTL active. Returns error if a problem occurs while trying to persist on cache.
func (r *redisPersistence) Lock(key string) custom_error.BaseErrorAdapter {
	return nil
}

// Unlock will remove the key from cache. Returns error if a problem occurs while trying to delete from cache.
func (r *redisPersistence) Unlock(key string) custom_error.BaseErrorAdapter {
	return nil
}

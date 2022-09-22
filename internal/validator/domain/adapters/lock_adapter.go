package adapters

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

type LockPersistenceAdapter interface {
	// Lock will set the key on cache with TTL active. Returns error if a problem occurs while trying to persist on cache.
	Lock(key string) custom_error.BaseErrorAdapter

	// Unlock will remove the key from cache. Returns error if a problem occurs while trying to delete from cache.
	Unlock(key string) custom_error.BaseErrorAdapter
}

package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type redisPersistence struct {
	LockCounter   int
	LockError     error
	UnlockCounter int
	UnlockError   error
	lock          map[string]string
}

func RedisPersistence() *redisPersistence {
	return &redisPersistence{
		lock: make(map[string]string),
	}
}

func (r *redisPersistence) Lock(key string) custom_error.BaseErrorAdapter {
	r.LockCounter++

	if r.LockError != nil {
		return exceptions.RedisPersistenceLockError(r.LockError, "Lock error")
	}

	if _, isLocked := r.lock[key]; isLocked {
		return exceptions.RedisPersistenceLockError(r.LockError, "key is already locked")
	}

	r.lock[key] = key

	return nil
}

func (r *redisPersistence) Unlock(key string) custom_error.BaseErrorAdapter {
	r.UnlockCounter++

	if r.UnlockError != nil {
		return exceptions.RedisPersistenceUnlockError(r.UnlockError, "Unlock error")
	}

	if _, isLocked := r.lock[key]; !isLocked {
		return exceptions.RedisPersistenceUnlockError(r.LockError, "key is already unlocked")
	}

	delete(r.lock, key)

	return nil
}

func (r *redisPersistence) IsLocked(key string) bool {
	_, isLocked := r.lock[key]
	return isLocked
}

func (r *redisPersistence) Reset() {
	r.LockCounter = 0
	r.LockError = nil
	r.UnlockCounter = 0
	r.UnlockError = nil
	r.lock = make(map[string]string)
}

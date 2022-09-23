package persistence

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/persistence"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	redisPersistence adapters2.LockPersistenceAdapter
	loggerM          = mocks.Logger()
	redis            = mocks.RedisServer()
)

var (
	key string
)

func redisPersistenceSetup() {
	config.LoadTestEnv()

	loggerM.Reset()
	redis.Reset()

	redisPersistence = persistence.RedisPersistence(loggerM, redis)

	key = uuid.NewString()
}

func teardown() {
	redis.Teardown()
}

func TestRedisLockSuccess(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	err := redisPersistence.Lock(key)

	keyLocked, _ := redis.Get(properties.Properties().Cache.KeyPrefix + key)

	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, keyLocked, key)
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 2, loggerM.InfoCallCounter)
	assert.Equal(t, 0, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockKeyAlreadyLockedFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)

	err := redisPersistence.Lock(key)

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "Key is already locked", err.Error())
	assert.Equal(t, "Key is already locked", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockKeyAlreadyLocked2Failure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	err := redisPersistence.Lock(key)

	assert.Nil(t, err, "Should be nil")

	err = redisPersistence.Lock(key)

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "Key is already locked", err.Error())
	assert.Equal(t, "Key is already locked", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 2, redis.OpenCounter)
	assert.Equal(t, 2, redis.CloseCounter)
	assert.Equal(t, 3, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockOpenConnectionFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.OpenError = errors.New("open conn error")

	err := redisPersistence.Lock(key)

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "open conn error", err.Error())
	assert.Equal(t, "Error while trying to open redis connection", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 0, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockGetFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.GetError = errors.New("get error")

	err := redisPersistence.Lock("error")

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "get error", err.Error())
	assert.Equal(t, "Error while trying to get redis key", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockSetFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.SetError = errors.New("set error")

	err := redisPersistence.Lock("error")

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "set error", err.Error())
	assert.Equal(t, "Error while trying to set redis key", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisLockCloseSuccess(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.CloseError = errors.New("close error")

	err := redisPersistence.Lock("error")

	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 2, loggerM.InfoCallCounter)
	assert.Equal(t, 0, loggerM.ErrorCallCounter)
	assert.Equal(t, 1, loggerM.WarningCallCounter)
}

func TestRedisUnlockSuccess(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)

	err := redisPersistence.Unlock(key)

	keyLocked, _ := redis.Get(properties.Properties().Cache.KeyPrefix + key)

	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, keyLocked, "")
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 2, loggerM.InfoCallCounter)
	assert.Equal(t, 0, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisUnlockKeyNotLockedSuccess(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	err := redisPersistence.Unlock(key)

	keyLocked, _ := redis.Get(properties.Properties().Cache.KeyPrefix + key)

	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, keyLocked, "")
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 2, loggerM.InfoCallCounter)
	assert.Equal(t, 0, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisUnlockOpenFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.OpenError = errors.New("open conn error")

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)
	err := redisPersistence.Unlock(key)

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "open conn error", err.Error())
	assert.Equal(t, "Error while trying to open redis connection", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 0, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisUnlockDelFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.DelError = errors.New("del error")

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)
	err := redisPersistence.Unlock("error")

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "del error", err.Error())
	assert.Equal(t, "Error while trying to delete redis key", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 0, loggerM.WarningCallCounter)
}

func TestRedisUnlockCloseSuccess(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.CloseError = errors.New("close error")

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)
	err := redisPersistence.Unlock(key)

	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 2, loggerM.InfoCallCounter)
	assert.Equal(t, 0, loggerM.ErrorCallCounter)
	assert.Equal(t, 1, loggerM.WarningCallCounter)
}

func TestRedisUnlockDelAndCloseErrorFailure(t *testing.T) {
	redisPersistenceSetup()
	defer teardown()

	redis.CloseError = errors.New("close error")
	redis.DelError = errors.New("del error")

	_ = redis.Set(properties.Properties().Cache.KeyPrefix+key, key)
	err := redisPersistence.Unlock("error")

	assert.NotNil(t, err, "Should not be nil")
	assert.Equal(t, "del error", err.Error())
	assert.Equal(t, "Error while trying to delete redis key", err.InternalError())
	assert.Equal(t, "Error while using Redis cache", err.Description())
	assert.Equal(t, 1, redis.OpenCounter)
	assert.Equal(t, 1, redis.CloseCounter)
	assert.Equal(t, 1, loggerM.InfoCallCounter)
	assert.Equal(t, 1, loggerM.ErrorCallCounter)
	assert.Equal(t, 1, loggerM.WarningCallCounter)
}

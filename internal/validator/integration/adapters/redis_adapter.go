package adapters

import "github.com/go-redis/redis/v8"

// RedisAdapter class for redis repository connection
type RedisAdapter interface {
	// Open connection to redis db
	Open() (*redis.Client, error)

	// Close connection to redis db
	Close() error
}

package mocks

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/go-redis/redis/v8"
	"strings"
)

type redisServer struct {
	server       *miniredis.Miniredis
	client       adapters.RedisAdapter
	OpenCounter  int
	OpenError    error
	CloseCounter int
	CloseError   error
	SetError     error
	GetError     error
	DelError     error
}

type redisTestHook struct {
	target string
	client *redis.Client
	match  string
	action func(*redis.Client) error
}

func RedisServer() *redisServer {
	redisServer := &redisServer{
		server: nil,
		client: nil,
	}
	redisServer.createClient()
	return redisServer
}

func (r *redisServer) createClient() adapters.RedisAdapter {
	if r.client != nil && r.server != nil {
		return r
	}

	config.LoadTestEnv()

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	secretsManagerService := SecretsManagerService()
	secretsManagerService.SetSecret(properties.Properties().Aws.SecretsManager.CacheSecretName, &dto.RedisSecrets{
		Address:    server.Addr(),
		Password:   "",
		User:       "",
		DatabaseId: 0,
	})

	client := config.RedisClient(secretsManagerService)

	r.server = server
	r.client = client

	return r
}

func (r *redisServer) Open() (*redis.Client, error) {
	r.OpenCounter++
	if r.OpenError != nil {
		return nil, r.OpenError
	}

	client, err := r.client.Open()
	if r.SetError != nil {
		client.AddHook(&redisTestHook{
			target: "before set",
			client: client,
			match:  "error",
			action: func(_ *redis.Client) error {
				return fmt.Errorf("set error")
			},
		})
	}
	if r.GetError != nil {
		client.AddHook(&redisTestHook{
			target: "before get",
			client: client,
			match:  "error",
			action: func(_ *redis.Client) error {
				return fmt.Errorf("get error")
			},
		})
	}
	if r.DelError != nil {
		client.AddHook(&redisTestHook{
			target: "before del",
			client: client,
			match:  "error",
			action: func(_ *redis.Client) error {
				return fmt.Errorf("del error")
			},
		})
	}
	return client, err
}

func (r *redisServer) Get(key string) (string, error) {
	redisClient, _ := r.client.Open()
	value, _ := redisClient.Get(redisClient.Context(), key).Result()
	err := redisClient.Close()

	return value, err
}

func (r *redisServer) Set(key string, value string) error {
	redisClient, _ := r.client.Open()
	_, _ = redisClient.Set(redisClient.Context(), key, value, 0).Result()
	err := redisClient.Close()

	return err
}

func (r *redisServer) Close() error {
	r.CloseCounter++
	if r.CloseError != nil {
		return r.CloseError
	}
	return r.client.Close()
}

func (r *redisServer) Teardown() {
	r.server.Close()
}

func (r *redisServer) Reset() {
	r.DelError = nil
	r.SetError = nil
	r.GetError = nil
	r.OpenCounter = 0
	r.OpenError = nil
	r.CloseCounter = 0
	r.CloseError = nil
	r.server = nil
	r.client = nil
	r.createClient()
}

func (r *redisTestHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if r.target == "before set" && strings.Contains(cmd.String(), r.match) && strings.Contains(cmd.String(), "set") {
		return nil, r.action(r.client)
	}
	if r.target == "before get" && strings.Contains(cmd.String(), r.match) && strings.Contains(cmd.String(), "get") {
		return nil, r.action(r.client)
	}
	if r.target == "before del" && strings.Contains(cmd.String(), r.match) && strings.Contains(cmd.String(), "del") {
		return nil, r.action(r.client)
	}
	return ctx, nil
}

func (r *redisTestHook) AfterProcess(context.Context, redis.Cmder) error {
	return nil
}

func (r *redisTestHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *redisTestHook) AfterProcessPipeline(context.Context, []redis.Cmder) error {
	return nil
}

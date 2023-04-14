package Utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	Host     string
	Port     string
	Password string
	Db       int
	Context  context.Context
}

func NewRedisClient(host string, port string, password string, db int) *RedisClient {
	return &RedisClient{
		Host:     host,
		Port:     port,
		Password: password,
		Db:       db,
		Context:  context.Background(),
	}
}

func (redisCli *RedisClient) Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisCli.Host + ":" + redisCli.Port,
		Password: redisCli.Password,
		DB:       redisCli.Db,
	})
	return client
}

func (redisCli *RedisClient) Set(key string, value string) error {
	client := redisCli.Connect()
	err := client.Set(redisCli.Context, key, value, 120*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (redisCli *RedisClient) Get(key string) (string, error) {
	client := redisCli.Connect()
	val, err := client.Get(redisCli.Context, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

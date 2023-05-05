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

func (redisCli *RedisClient) LRange(key string, start int64, stop int64) ([]string, error) {
	client := redisCli.Connect()
	val, err := client.LRange(redisCli.Context, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (redisCli *RedisClient) RPush(key string, value string) error {
	client := redisCli.Connect()
	err := client.RPush(redisCli.Context, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (redisCli *RedisClient) Del(key string) error {
	client := redisCli.Connect()
	err := client.Del(redisCli.Context, key).Err()
	if err != nil {
		return err
	}
	return nil
}

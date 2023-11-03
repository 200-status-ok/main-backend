package Utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/dtos"
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

func (redisCli *RedisClient) Del(key string) error {
	client := redisCli.Connect()
	err := client.Del(redisCli.Context, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (redisCli *RedisClient) PublishMessageToUserChannel(channel string, message dtos.Message) error {
	client := redisCli.Connect()
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = client.Publish(redisCli.Context, channel, jsonMessage).Err()
	if err != nil {
		return err
	}
	return nil
}

func (redisCli *RedisClient) SubscribeToUserChannel(channel string, messageChan chan<- dtos.Message) {
	client := redisCli.Connect()
	pubsub := client.Subscribe(redisCli.Context, channel)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {

		}
	}(pubsub)

	for msg := range pubsub.Channel() {
		var message dtos.Message
		err := json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messageChan <- message
	}
}

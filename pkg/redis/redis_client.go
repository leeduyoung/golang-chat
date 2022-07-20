package redis

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/LUXROBO/server/libs/env"
	"github.com/go-redis/redis/v8"
)

const expirationTime = 0

var (
	once          sync.Once
	redisInstance *redisClientService
)

type redisClientService struct {
	ctx         context.Context
	redisClient *redis.Client
}

type IRedisClientService interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	SIsMember(key string, value interface{}) (bool, error)
	Publish(msgToPub MsgToPub)
	Subscribe(msgToSub MsgToSub) *redis.PubSub
}

type MsgToPub struct {
	Channel string
	Message string
}

type MsgToSub struct {
	Channels []string
}

func GetInstance() IRedisClientService {
	if redisInstance == nil {
		host := env.Instance.ChatRedisHost
		port := env.Instance.ChatRedisPort
		password := env.Instance.ChatRedisPassword

		once.Do(func() {
			redisInstance = &redisClientService{
				redisClient: redis.NewClient(&redis.Options{
					Addr:     host + ":" + port,
					Password: password,
				}),
				ctx: context.Background(),
			}
		})
	}

	return redisInstance
}

func (rcs *redisClientService) Set(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	return rcs.redisClient.Set(rcs.ctx, key, string(bytes), expirationTime).Err()
}

func (rcs *redisClientService) Get(key string) (string, error) {
	return rcs.redisClient.Get(rcs.ctx, key).Result()
}

func (rcs *redisClientService) SIsMember(key string, value interface{}) (bool, error) {
	cmd := rcs.redisClient.SIsMember(rcs.ctx, key, value)
	return cmd.Result()
}

func (rcs *redisClientService) Publish(msgToPub MsgToPub) {
	if err := rcs.redisClient.Publish(rcs.ctx, msgToPub.Channel, msgToPub.Message).Err(); err != nil {
		panic(err)
	}
}

func (rcs *redisClientService) Subscribe(msgToSub MsgToSub) *redis.PubSub {
	return rcs.redisClient.Subscribe(rcs.ctx, msgToSub.Channels...)
}

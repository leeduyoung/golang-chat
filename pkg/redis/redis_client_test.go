package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/LUXROBO/server/libs/env"
	"github.com/LUXROBO/server/libs/logger"
	"github.com/LUXROBO/server/services/chat/pkg/event"
)

const serviceName = "chat"

func init() {
	env := env.Initialize(serviceName)
	logger.Initialize(serviceName, env.Mode)
}

func TestRedisSetGet(t *testing.T) {
	const (
		roomName = "test-room-1"
		userID   = "test-user-1"
		userID2  = "test-user-1"
	)

	users := []string{userID, userID2}

	// SET
	err := GetInstance().Set(roomName, users)
	if err != nil {
		t.Error(err)
		return
	}

	// GET
	val, err := GetInstance().Get(roomName)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("[Success] RedisSetGet - val: ", val)
}

func TestRedisPubSub(t *testing.T) {
	redisPubSubChannel := "test-channel"

	go t.Run("subscribe", func(t *testing.T) {
		msgToSub := MsgToSub{
			Channels: []string{redisPubSubChannel},
		}

		subscriber := GetInstance().Subscribe(msgToSub)

		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				panic(err)
			}

			logger.Info("msg: ", msg)
		}
	})

	time.Sleep(time.Second)

	t.Run("publish", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			msg := event.ChatEvent{
				UserID:      fmt.Sprintf("user[%d]", i+1),
				RoomName:    fmt.Sprintf("room[%d]", i+1),
				Nickname:    fmt.Sprintf("nickname[%d]", i+1),
				Message:     "test",
				MessageType: event.MessageTypeEnter,
			}
			bytes, _ := json.Marshal(msg)

			msgToPub := MsgToPub{
				Channel: redisPubSubChannel,
				Message: string(bytes),
			}

			GetInstance().Publish(msgToPub)
		}
	})
}

package service

import (
	"context"
	"encoding/json"

	"github.com/LUXROBO/server/libs/logger"
	"github.com/LUXROBO/server/services/chat/pkg/event"
	"github.com/LUXROBO/server/services/chat/pkg/helper/array"
	redisUtils "github.com/LUXROBO/server/services/chat/pkg/helper/redis-utils"
	"github.com/LUXROBO/server/services/chat/pkg/redis"
)

// RoomMessage 특정방에 전달될 메시지와 타깃 목록 (유저 아이디)
type RoomMessage struct {
	targetIDs []string
	message   []byte
}

// Hub 소켓서버에 연결된 유저들을 저장하고 유저에게 메시지를 전달하는 구조체
type Hub struct {
	enter          chan *Client
	leave          chan *Client
	roomMessage    chan RoomMessage
	userIDToClient map[string]*Client
}

// NewHub Hub 구조체 생성
func NewHub() *Hub {
	return &Hub{
		enter:          make(chan *Client),
		leave:          make(chan *Client),
		roomMessage:    make(chan RoomMessage),
		userIDToClient: make(map[string]*Client),
	}
}

// Run 소켓을 통해 메시지를 전달
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.enter:
			h.userIDToClient[client.userID] = client
		case client := <-h.leave:
			if _, ok := h.userIDToClient[client.userID]; ok {
				delete(h.userIDToClient, client.userID)
				close(client.send)
			}
		case roomMessage := <-h.roomMessage:
			for _, userID := range roomMessage.targetIDs {
				if val, ok := h.userIDToClient[userID]; ok {
					val.send <- roomMessage.message
				}
			}
		}
	}
}

// Subscribe redis의 채팅 메시지를 수신하는 함수
func (h *Hub) Subscribe() {
	ctx := context.Background()
	subscriber := redis.GetInstance().Subscribe(redis.MsgToSub{
		Channels: []string{event.ChatChannel},
	})

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		logger.Info("[subscribe] msg: ", msg)

		data, err := parseChatMessage(msg.Payload)
		if err != nil {
			logger.Error("[subscribe] parseChatMessage error: ", err)
			continue
		}

		users, err := redisUtils.LoadUserList(data.RoomName)
		if err != nil {
			logger.Error("[subscribe] loadUserList error: ", err)
			continue
		}

		h.roomMessage <- RoomMessage{
			targetIDs: array.RemoveDuplcateItem(users),
			message:   []byte(msg.Payload),
		}
	}
}

func parseChatMessage(payload string) (event.ChatEvent, error) {
	data := event.ChatEvent{}
	err := json.Unmarshal([]byte(payload), &data)
	return data, err
}

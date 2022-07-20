package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LUXROBO/server/libs/logger"
	"github.com/LUXROBO/server/services/chat/pkg/event"
	"github.com/LUXROBO/server/services/chat/pkg/helper/array"
	redisUtils "github.com/LUXROBO/server/services/chat/pkg/helper/redis-utils"
	"github.com/LUXROBO/server/services/chat/pkg/redis"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client websocket 연결과 hub 사이를 연결하는 중개자
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

// readPump 웹소켓으로 전달받은 데이터를 redis로 publish 처리
func (c *Client) readPump() {
	defer func() {
		c.hub.leave <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		msg, err := parseSocketMessage(message)
		if err != nil {
			logger.Error("[readPump] filaed parse socket message. message: ", string(message))
		}

		users, err := redisUtils.LoadUserList(msg.RoomName)
		if err != nil {
			logger.Warn("[readPump] loadUserList error: ", err)
		}

		if len(users) == 0 || !array.Contains(users, msg.UserID) {
			users = append(users, msg.UserID)
			redis.GetInstance().Set(msg.RoomName, users)
		}

		redis.GetInstance().Publish(redis.MsgToPub{
			Channel: event.ChatChannel,
			Message: string(message),
		})
	}
}

// writePump hub를 통해서 전달받은 메시지를 웹소켓으로 전달
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServerWs 웹소켓 요청을 처리하는 함수
func ServerWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.enter <- client

	go client.writePump()
	go client.readPump()
}

func parseSocketMessage(message []byte) (*event.ChatEvent, error) {
	msg := &event.ChatEvent{}
	err := json.Unmarshal(message, msg)
	return msg, err
}

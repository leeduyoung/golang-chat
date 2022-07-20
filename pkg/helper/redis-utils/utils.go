package redisUtils

import (
	"encoding/json"

	"github.com/LUXROBO/server/services/chat/pkg/redis"
)

func LoadUserList(roomName string) ([]string, error) {
	users := []string{}
	response, err := redis.GetInstance().Get(roomName)

	if err != nil {
		return users, err
	}

	err = json.Unmarshal([]byte(response), &users)
	return users, err
}

package event

type MessageType string

const (
	MessageTypeEnter MessageType = "ENTER"
	MessageTypeLeave MessageType = "LEAVE"
	MessageTypeChat  MessageType = "CHAT"
	ChatChannel                  = "atc-chat-channel"
)

type ChatEvent struct {
	UserID      string      `json:"user_id"`
	Nickname    string      `json:"nickname"`
	RoomName    string      `json:"roomName"`
	Message     string      `json:"message"`
	MessageType MessageType `json:"messageType"`
}

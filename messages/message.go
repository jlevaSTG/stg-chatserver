package messages

import (
	"time"
)

var (
	Connected        = "connect"
	Disconnected     = "disconnect"
	ConnectionStatus = "connection_status"
	TextChatMessage  = "text_chat_message"
)

type Message struct {
	MessageType string      `json:"message_type"`
	CreatedAt   time.Time   `json:"created_at"`
	Payload     interface{} `json:"payload"`
}

type UserConnectionStatus struct {
	ClientID        string `json:"client_id"`
	ConnectedStatus string `json:"connected_status"`
}

type ChatMessage struct {
	ChatId      string `json:"chat_id" bson:"chat_id"`
	SentBy      string `json:"sent_by" bson:"sent_by"`
	Message     string `json:"message" bson:"message"`
	ResourceUrl string `json:"resource_url" bson:"resource_url"`
}

func NewMessage(messageType string, payload interface{}) Message {
	return Message{
		MessageType: messageType,
		CreatedAt:   time.Now(),
		Payload:     payload,
	}
}

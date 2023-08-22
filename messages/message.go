package messages

import (
	"time"
)

var (
	Connected        = "connect"
	Disconnected     = "disconnect"
	ConnectionStatus = "connection_status"
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

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stg-go-websocket-server/types"
	"time"
)

type MessageType string

var (
	TextMessageType MessageType = "text_message"
)

type ChatSession struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"-"`
	ChatId       string              `bson:"chat_id" json:"chat_id"`
	CreatedAt    time.Time           `bson:"created_at" json:"created_at"`
	CreatedBy    string              `bson:"created_by" json:"created_by"`
	Active       bool                `bson:"active" json:"active"`
	Participants []types.Participant `bson:"participants" json:"participants"`
	Messages     []ChatMessage       `bson:"messages" json:"messages"`
}

func NewChatSession(chatId string, createdBY string, participants []types.Participant, m string) ChatSession {
	for i := range participants {
		participants[i].Active = true
		participants[i].AddedBy = createdBY
		participants[i].JoinedAt = time.Now()
	}

	return ChatSession{
		ChatId:       chatId,
		CreatedBy:    createdBY,
		CreatedAt:    time.Now(),
		Active:       true,
		Participants: participants,
		Messages:     []ChatMessage{NewTextChatMessage(chatId, createdBY, m)},
	}
}

type ChatMessage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	MessageType MessageType        `json:"message_type" bson:"message_type"`
	ChatID      string             `json:"chat_id" bson:"chat_id"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	Message     string             `json:"message" bson:"message"`
	ResourceUrl string             `json:"resource_url" bson:"resource_url"`
}

func NewTextChatMessage(chatID string, createdBy string, message string) ChatMessage {
	return ChatMessage{
		CreatedAt:   time.Now(),
		MessageType: TextMessageType,
		ChatID:      chatID,
		CreatedBy:   createdBy,
		Message:     message,
	}
}

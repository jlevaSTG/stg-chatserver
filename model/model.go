package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stg-go-websocket-server/types"
	"time"
)

type ChatSession struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"-"`
	ChatId       string              `bson:"chat_id"`
	CreatedAt    time.Time           `bson:"created_at"`
	CreatedBy    string              `bson:"created_by"`
	Active       bool                `bson:"active"`
	Participants []types.Participant `bson:"participants"`
}

func NewChatSession(chatId string, createdBY string, participants []types.Participant) ChatSession {
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
	}
}

type ChatMessage struct {
	MessageType string `json:"message_type" bson:"message_type"`
	ChatId      string `json:"chat_id" bson:"chat_id"`
	SentBy      string `json:"sent_by" bson:"sent_by"`
	Message     string `json:"message" bson:"message"`
	ResourceUrl string `json:"resource_url" bson:"resource_url"`
}

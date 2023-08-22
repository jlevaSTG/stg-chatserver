package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stg-go-websocket-server/types"
	"time"
)

type ChatSession struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty"`
	ChatId       string              `bson:"chat_id"`
	CreatedAt    int64               `bson:"created_at"`
	CreatedBy    string              `bson:"created_by"`
	Active       bool                `bson:"active"`
	Participants []types.Participant `bson:"participants"`
}

func NewChatSession(chatId string, createdBY string, p []types.Participant) ChatSession {
	return ChatSession{
		ChatId:       chatId,
		CreatedBy:    createdBY,
		CreatedAt:    time.Now().Unix(),
		Active:       true,
		Participants: p,
	}
}

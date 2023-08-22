package ws

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"stg-go-websocket-server/model"
	"stg-go-websocket-server/types"
	"time"
)

type CommandHandler func(cmd Command) error

const (
	InitChatCommandType = "init-chat"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("stg-chat").Collection(collectionName)
	return collection
}

type Command struct {
	CommandType      string          `json:"cmd-type"`
	CommandTimeStamp string          `json:"cmd-time-stamp"`
	CommandIssuer    string          `json:"command-issuer"`
	Payload          json.RawMessage `json:"payload"`
}

type InitChatCommand struct {
	ClientID     string              `json:"id"`
	CreatedAt    time.Time           `json:"created_at"`
	CreatedBy    string              `json:"created_by"`
	Participants []types.Participant `json:"participants"`
	Message      string              `json:"message"`
}

func (cmd *InitChatCommand) participantsIds() []string {
	ids := make([]string, len(cmd.Participants))
	for _, p := range cmd.Participants {
		ids = append(ids, p.ID)
	}
	return ids
}

func SetUpCommandHandlers(m *Manager) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[InitChatCommandType] = InitChatCommandHandlers(m)
	return handlers
}

func InitChatCommandHandlers(m *Manager) CommandHandler {
	return func(cmd Command) error {
		var initCmd InitChatCommand
		if err := json.Unmarshal(cmd.Payload, &initCmd); err != nil {
			return err
		}

		dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection := GetCollection(m.mongoClient, "chat_sessions")
		// Sort the participants to ensure consistent query
		ids := initCmd.participantsIds()
		sort.Strings(ids)
		// Check if a chat session with these participants already exists
		var existingSession model.ChatSession
		err := collection.FindOne(dbCtx, bson.M{"participants": ids}).Decode(&existingSession)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				// No existing session, go ahead to create new one
				// Your code for creating new session
				s := model.NewChatSession(m.IdCreator(), initCmd.CreatedBy, initCmd.Participants)
				_, err := collection.InsertOne(dbCtx, s)
				if err != nil {
					log.Err(err)
					return err
				}

			}
		}
		// Existing session found send messages

		return nil

	}
}

//func (m *Manager) startChatSession(ctx *gin.Context) {
//	var request chat.InitSession
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	collection := GetCollection(m.mongoClient, "chat_sessions")
//	// Sort the participants to ensure consistent query
//	sort.Strings(request.Participants)
//	// Check if a chat session with these participants already exists
//	var existingSession chat.InitSession
//	err := collection.FindOne(dbCtx, bson.M{"participants": request.Participants}).Decode(&existingSession)
//
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			// No existing session, go ahead to create new one
//			// Your code for creating new session
//			ctx.JSON(http.StatusOK, gin.H{"message": "New chat session started"})
//			return
//		}
//		// Other errors
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	// Existing session found
//	ctx.JSON(http.StatusConflict, gin.H{"message": "Chat session already exists", "session": existingSession})
//
//}

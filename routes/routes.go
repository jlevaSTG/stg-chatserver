package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sort"
	"stg-go-websocket-server/model"
	"stg-go-websocket-server/types"
	"stg-go-websocket-server/ws"
	"time"
)

func SetupApiRoutes(route string, r *gin.Engine, m *ws.Manager) {
	api := r.Group(route)
	api.POST("/initChat", handleChatInit(m))
	api.POST("/sendMsg", handleSendMsg(m))
	api.GET("/retrieve/:userID", retrieveChat(m))
}

type RetrieveChat struct {
	ChatSession []ChatSession `json:"chat_session"`
}

type ChatSession struct {
	ChatId       string              `json:"chat_id"`
	CreatedAt    time.Time           `json:"created_at"`
	CreatedBy    string              `json:"created_by"`
	Active       bool                `json:"active"`
	Participants []types.Participant `json:"participants"`
	ChatMessages []model.ChatMessage `json:"chat_messages"`
}

func retrieveChat(m *ws.Manager) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		userID := ginCtx.Param("userID")

		// Get MongoDB collection
		sessions := ws.GetCollection(m.MongoClient, "chat_sessions")
		dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Filter to find all sessions where userID is a participant
		filter := bson.M{"participants": bson.M{"$elemMatch": bson.M{"id": userID}}}
		cursor, err := sessions.Find(dbCtx, filter)

		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
			return
		}

		var chatSessions []model.ChatSession
		if err = cursor.All(dbCtx, &chatSessions); err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode sessions"})
			return
		}

		ids := make([]string, 0)
		for _, s := range chatSessions {
			ids = append(ids, s.ChatId)
		}
		chatMessages := ws.GetCollection(m.MongoClient, "chat_messages")
		filter = bson.M{"chat_id": bson.M{"$in": ids}}
		cursor, err = chatMessages.Find(dbCtx, filter)

		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat messages"})
			return
		}

		var fetchedChatMessages []model.ChatMessage // Replace YourChatMessageModel with the actual model type
		if err = cursor.All(dbCtx, &fetchedChatMessages); err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode chat messages"})
			return
		}

		chatSessionReturn := make([]ChatSession, 0)
		for _, s := range chatSessions {
			chatSessionReturn = append(chatSessionReturn, ChatSession{
				ChatId:       s.ChatId,
				CreatedAt:    s.CreatedAt,
				CreatedBy:    s.CreatedBy,
				Active:       s.Active,
				Participants: s.Participants,
				ChatMessages: getMessages(fetchedChatMessages, s.ChatId),
			})
		}

		// Use chatSessions...
		ginCtx.JSON(http.StatusOK, gin.H{"chatSessions": chatSessionReturn})
	}
}

func getMessages(messages []model.ChatMessage, chatId string) []model.ChatMessage {
	msg := make([]model.ChatMessage, 0)
	for _, m := range messages {
		if m.ChatID == chatId {
			msg = append(msg, m)
		}
	}
	return msg
}

func handleSendMsg(m *ws.Manager) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var msgCmd ws.TextMessageCommand
		if err := bindJSON(ginCtx, &msgCmd); err != nil {
			return
		}

		existingSession, err := findExistingSession(m.MongoClient, msgCmd.ChatID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "no chat session with id " + msgCmd.ChatID + " found"})
				return
			}
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = createChatMessage(m.MongoClient, msgCmd, existingSession)
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sendCommand(m, msgCmd, existingSession)
		ginCtx.JSON(http.StatusOK, gin.H{"chat": msgCmd})
	}
}

func bindJSON[T any](ginCtx *gin.Context, cmd T) error {
	if err := ginCtx.ShouldBindJSON(cmd); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

func findExistingSession(mongoClient *mongo.Client, chatID string) (*model.ChatSession, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionCollection := ws.GetCollection(mongoClient, "chat_sessions")
	filter := bson.M{"chat_id": chatID}

	var existingSession model.ChatSession
	err := sessionCollection.FindOne(dbCtx, filter).Decode(&existingSession)
	if err != nil {
		return nil, err
	}

	return &existingSession, nil
}

func createChatMessage(mongoClient *mongo.Client, msgCmd ws.TextMessageCommand, existingSession *model.ChatSession) (*model.ChatMessage, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &model.ChatMessage{
		MessageType: model.TextMessageType,
		ChatID:      msgCmd.ChatID,
		CreatedBy:   msgCmd.CreatedBy,
		Message:     msgCmd.Message,
	}

	chatCollection := ws.GetCollection(mongoClient, "chat_messages")
	_, err := chatCollection.InsertOne(dbCtx, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func sendCommand(m *ws.Manager, msgCmd ws.TextMessageCommand, existingSession *model.ChatSession) {
	cmd := ws.NewCommand(ws.TextMessageCommandType, msgCmd.CreatedBy, ws.TextMessageCommand{
		CreatedBy:       msgCmd.CreatedBy,
		ChatID:          msgCmd.ChatID,
		ParticipantsIds: existingSession.Participants,
		Message:         msgCmd.Message,
	})
	m.CommandStream <- cmd
}

func handleChatInit(m *ws.Manager) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var initCmd ws.InitChatCommand
		if err := ginCtx.ShouldBindJSON(&initCmd); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Info().Msgf("init chat for cmd: %v", initCmd)

		dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection := ws.GetCollection(m.MongoClient, "chat_sessions")

		// Sort the participants to ensure consistent query
		ids := initCmd.ParticipantsIds()
		sort.Strings(ids)

		filter := bson.M{
			"$and": []bson.M{
				{"participants.id": bson.M{"$all": ids}},
				{"participants": bson.M{"$size": len(ids)}},
			},
		}

		var existingSession model.ChatSession
		err := collection.FindOne(dbCtx, filter).Decode(&existingSession)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err == nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "chat between participants already exists"})
			return
		}
		log.Info().Msgf("no session found creating new chat: %v", initCmd)

		s := model.NewChatSession(m.IdCreator(), initCmd.CreatedBy, initCmd.Participants)
		if _, err := collection.InsertOne(dbCtx, s); err != nil {
			log.Err(err).Msg("Failed to insert new chat session")
		}

		cmd := ws.NewCommand(ws.TextMessageCommandType, s.CreatedBy, ws.TextMessageCommand{
			CreatedBy:       s.CreatedBy,
			ChatID:          s.ChatId,
			ParticipantsIds: s.Participants,
			Message:         initCmd.Message,
		})

		m.CommandStream <- cmd
		ginCtx.JSON(http.StatusCreated, gin.H{"chat": s})
	}
}

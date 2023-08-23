package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sort"
	"stg-go-websocket-server/model"
	"stg-go-websocket-server/ws"
	"time"
)

func SetupApiRoutes(route string, r *gin.Engine, m *ws.Manager) {
	api := r.Group(route)
	api.POST("/initChat", handleChatInit(m))
	api.POST("/sendMsg", handleSendMsg(m))
}

func handleSendMsg(m *ws.Manager) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var msgCmd ws.TextMessageCommand
		if err := ginCtx.ShouldBindJSON(&msgCmd); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Info().Msgf("send msg for chat: %v", msgCmd)

		dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		collection := ws.GetCollection(m.MongoClient, "chat_sessions")
		fmt.Println("Debug: ChatID is ", msgCmd.ChatID) // Debug
		filter := bson.M{"chat_id": msgCmd.ChatID}
		fmt.Println("Debug: Filter is ", filter)
		var existingSession model.ChatSession
		err := collection.FindOne(dbCtx, filter).Decode(&existingSession)

		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(" no chat found with id %s", msgCmd.ChatID)})
			return
		}
		log.Info().Msgf("no session found creating new chat: %v", msgCmd)
		cmd := ws.NewCommand(ws.TextMessageCommandType, msgCmd.CreatedBy, ws.TextMessageCommand{
			CreatedBy:       msgCmd.CreatedBy,
			ChatID:          msgCmd.ChatID,
			ParticipantsIds: existingSession.Participants,
			Message:         msgCmd.Message,
		})

		m.CommandStream <- cmd
		ginCtx.JSON(http.StatusOK, gin.H{"chat": msgCmd})

	}
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

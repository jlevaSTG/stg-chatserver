package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"stg-go-websocket-server/messages"
	"sync"
	"time"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
	// pongWait is how long we will await a pong response from client
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type IdCreator = func() string

type Manager struct {
	sync.RWMutex
	CurrentClientCount int64
	Clients            map[string]*Client
	handlers           map[string]CommandHandler
	CommandStream      chan Command
	IdCreator          IdCreator
	MongoClient        *mongo.Client
}

func NewManager(mongoClient *mongo.Client) *Manager {
	m := &Manager{
		MongoClient:   mongoClient,
		Clients:       make(map[string]*Client),
		CommandStream: make(chan Command, 50),
		IdCreator:     func() string { return uuid.New().String() },
	}
	m.handlers = SetUpCommandHandlers(m)
	m.startCmdLoop()
	return m
}

func (m *Manager) HandleWS(ctx *gin.Context) {
	userID, exists := ctx.GetQuery("userId")
	if !exists {
		ctx.JSON(http.StatusBadRequest, errors.New("user ClientID name must be sent"))
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade to WebSocket"})
		return
	}

	newClient := NewClient(userID, conn)
	m.Lock()
	m.Clients[newClient.Id] = newClient
	m.Unlock()
	m.CurrentClientCount++
	log.Info().Msgf("manager current client count: %d, Clients %v", m.CurrentClientCount, m.Clients)

	newClient.StartWriteLoop()
	newClient.StartReadLoop(m.CommandStream)

	msg := messages.Message{
		MessageType: messages.ConnectionStatus,
		CreatedAt:   time.Now(),
		Payload: messages.UserConnectionStatus{
			ClientID:        newClient.Id,
			ConnectedStatus: messages.Connected,
		},
	}

	newClient.Egress <- msg

}

func (m *Manager) startCmdLoop() {
	go func() {
		for cmd := range m.CommandStream {
			if handler, ok := m.handlers[cmd.CommandType]; ok {
				if err := handler(cmd); err != nil {
					log.Printf("Error executing command %v %v", err, cmd)
				}
			} else {
				log.Err(ErrEventNotSupported)
			}
		}
	}()
}

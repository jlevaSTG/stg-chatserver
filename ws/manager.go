package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
	currentClientCount int64
	clients            map[string]*Client
	handlers           map[string]CommandHandler
	commandStream      chan Command
	IdCreator          IdCreator
	mongoClient        *mongo.Client
}

func NewManager(mongoClient *mongo.Client) *Manager {
	m := &Manager{
		mongoClient:   mongoClient,
		clients:       make(map[string]*Client),
		commandStream: make(chan Command, 50),
		IdCreator:     uuid.New().String,
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
	m.clients[newClient.Id] = newClient
	m.Unlock()

	newClient.StartWriteLoop()
	newClient.StartReadLoop(m.commandStream)

	msg := messages.Message{
		MessageType: messages.ConnectionStatus,
		CreatedAt:   time.Now(),
		Payload: messages.UserConnected{
			ClientID:        newClient.Id,
			ConnectedStatus: messages.Connected,
		},
	}

	newClient.Egress <- msg

}

func (m *Manager) startCmdLoop() {
	go func() {
		for cmd := range m.commandStream {
			if handler, ok := m.handlers[cmd.CommandType]; ok {
				if err := handler(cmd); err != nil {
					log.Printf("Error executing command %v\n", err)
				}
			} else {
				log.Printf("Error executing command %v\n", ErrEventNotSupported)
			}
		}
	}()
}

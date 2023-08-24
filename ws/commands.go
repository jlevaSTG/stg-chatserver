package ws

import (
	"errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"stg-go-websocket-server/messages"
	"stg-go-websocket-server/types"
	"time"
)

type CommandHandler func(cmd Command) error

const (
	DisconnectClientCommandType = "disconnect-cmd"
	TextMessageCommandType      = "text-message-cmd"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("stg-chat").Collection(collectionName)
	return collection
}

type Command struct {
	CommandType   string      `json:"cmd-type"`
	CommandIssuer string      `json:"command-issuer"`
	Payload       interface{} `json:"payload"`
}

func NewCommand(commandType string, commandIssuer string, payload interface{}) Command {
	return Command{
		CommandType:   commandType,
		CommandIssuer: commandIssuer,
		Payload:       payload,
	}
}

type InitChatCommand struct {
	ClientID     string              `json:"id"`
	CreatedAt    time.Time           `json:"created_at"`
	CreatedBy    string              `json:"created_by"`
	Participants []types.Participant `json:"participants"`
	Message      string              `json:"message"`
}

type TextMessageCommand struct {
	CreatedBy       string              `json:"created_by"`
	ChatID          string              `json:"chat_id"`
	ParticipantsIds []types.Participant `json:"participants"`
	Message         string              `json:"message"`
}

type RetrieveChatCommand struct {
	USerID string `json:"userID"`
}

type RetrieveChatMessagesCommand struct {
	ChatIDS []string `json:"chat_ids"`
}

type DisconnectCommand struct {
	ClientID string `json:"id"`
}

func (cmd *InitChatCommand) ParticipantsIds() []string {
	ids := make([]string, len(cmd.Participants))
	for i, p := range cmd.Participants {
		ids[i] = p.ID
	}
	return ids
}

func SetUpCommandHandlers(m *Manager) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[DisconnectClientCommandType] = removeClientCommandHandlers(m)
	handlers[TextMessageCommandType] = handleTextMessageCommand(m)
	return handlers
}

func handleTextMessageCommand(m *Manager) CommandHandler {
	return func(cmd Command) error {
		switch payload := cmd.Payload.(type) {
		case TextMessageCommand:
			chatMsg := messages.ChatMessage{
				ChatId:  payload.ChatID,
				SentBy:  payload.CreatedBy,
				Message: payload.Message,
			}

			clientMsg := messages.NewMessage(messages.TextChatMessage, chatMsg)
			for _, p := range payload.ParticipantsIds {
				if p.ID != chatMsg.SentBy {
					c, ok := m.Clients[p.ID]
					if ok {
						c.Egress <- clientMsg
					}
				}
			}
		default:
			return errors.New("text command sent with wrong payload type")
		}
		return nil
	}
}

func removeClientCommandHandlers(m *Manager) CommandHandler {
	return func(cmd Command) error {
		switch payload := cmd.Payload.(type) {
		case DisconnectCommand:
			clientId := payload.ClientID
			log.Info().Msgf("removing client %s", clientId)

			m.Lock()
			client, ok := m.Clients[clientId]
			if ok {
				delete(m.Clients, clientId)
				err := client.Connection.Close()
				if err != nil {
					return err
				}
			}
			m.CurrentClientCount--
			m.Unlock()
			log.Info().Msgf("manager current client count: %d, Clients %v", m.CurrentClientCount, m.Clients)
		default:
			return errors.New("text command sent with wrong payload type")
		}
		return nil
	}

}

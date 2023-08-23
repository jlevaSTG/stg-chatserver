package ws

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"stg-go-websocket-server/messages"
	"stg-go-websocket-server/types"
	"time"
)

type CommandHandler func(cmd Command) error

const (
	DisconnectClientCommandType = "disconnect-cmd"
	InitChatCommandType         = "init-chat-cmd"
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
	handlers[InitChatCommandType] = InitChatCommandHandlers(m)
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
				c, ok := m.clients[p.ID]
				if ok {
					c.Egress <- clientMsg
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
		//var disconnect DisconnectCommand
		//if err := json.Unmarshal(cmd.Payload, &disconnect); err != nil {
		//	return err
		//}
		//
		//clientId := disconnect.ClientID
		//log.Info().Msgf("removing client %s\n", clientId)
		//
		//m.Lock()
		//client, ok := m.clients[clientId]
		//if ok {
		//	delete(m.clients, clientId)
		//	err := client.Connection.Close()
		//	if err != nil {
		//		return err
		//	}
		//}
		//m.currentClientCount--
		//m.Unlock()
		//log.Info().Msgf("manager current client count: %d, clients %v", m.currentClientCount, m.clients)
		return nil
	}
}

func InitChatCommandHandlers(m *Manager) CommandHandler {
	return func(cmd Command) error {
		//var initCmd InitChatCommand
		//if err := json.Unmarshal(cmd.Payload, &initCmd); err != nil {
		//	return err
		//}
		//log.Info().Msgf("init chat for cmd: %v", initCmd)
		//
		//chatMsg := messages.ChatMessage{
		//	ChatId:  s.ChatId,
		//	SentBy:  s.CreatedBy,
		//	Message: initCmd.Message,
		//}
		//
		//clientMsg := messages.NewMessage(messages.TextChatMessage, chatMsg)
		//
		//for _, p := range initCmd.Participants {
		//	c, ok := m.clients[p.ID]
		//	if ok {
		//		c.Egress <- clientMsg
		//	}
		//}

		return nil

	}
}

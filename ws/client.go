package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"stg-go-websocket-server/messages"
	"time"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	Id         string
	Joined     time.Time
	Connection *websocket.Conn
	Egress     chan messages.Message
}

func NewClient(id string, connection *websocket.Conn) *Client {
	return &Client{
		Id:         id,
		Joined:     time.Now(),
		Connection: connection,
		Egress:     make(chan messages.Message),
	}

}

func (c *Client) StartReadLoop(cmdStream chan Command) {
	go func() {
		defer func() {
			c.disconnect(cmdStream)
		}()

		if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Println(err)
			return
		}

		c.Connection.SetPongHandler(func(appData string) error {
			return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		})

		for {
			_, payload, err := c.Connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error reading message: %v", err)
				}
				break
			}
			var cmd Command
			if err := json.Unmarshal(payload, &cmd); err != nil {
				log.Printf("error marshalling message: %v", err)
			}
			cmdStream <- cmd
		}
	}()
}

func (c *Client) StartWriteLoop() {
	ticker := time.NewTicker(pingInterval)
	go func() {
		defer func() {
			ticker.Stop()
		}()
		for {
			select {
			case msg, ok := <-c.Egress:
				if !ok {
					if err := c.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
						log.Println("connection closed: ", err)
					}
					return
				}
				data, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
				}
				if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Println(err)
				}
				log.Println("sent message")
			case <-ticker.C:
				//fmt.Println("ping")
				if err := c.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.Println("write Msg: ", err)
					return
				}
			}
		}
	}()
}

func (c *Client) disconnect(cmdStream chan Command) {
	cmd := Command{CommandType: DisconnectClientCommandType, Payload: DisconnectCommand{ClientID: c.Id}}
	cmdStream <- cmd
}

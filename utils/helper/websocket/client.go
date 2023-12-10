package websocket

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	Message  chan *Message
}

type Message struct {
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *Client) WriteMessage(roomID string) {

	defer func() {
		c.Conn.Close()
	}()

	var messages []schema.Message

	configs.DB.Where("room_id = ?", roomID).Find(&messages)
	c.Conn.WriteJSON(messages)

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {

	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		message := schema.Messages{
			RoomID:  msg.RoomID,
			Message: msg.Content,
		}

		configs.DB.Create(&message)

		hub.Broadcast <- msg
	}
}

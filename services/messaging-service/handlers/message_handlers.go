package handlers

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/temuka-messaging-service/config"
	"github.com/temuka-messaging-service/models"
	"github.com/temuka-messaging-service/pb"
	"github.com/temuka-messaging-service/ws"
	"gorm.io/gorm"
)

type Client struct {
	Conn        *websocket.Conn
	Message     chan *models.Message
	Participant *models.Participant
	Hub         *ws.Hub
	DB          *gorm.DB
}

func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	db := config.GetDBInstance()

	newMessage := models.Message{
		ParticipantID: int(req.ParticipantId),
		Text:          req.Text,
	}

	if err := db.Create(&newMessage).Error; err != nil {
		return nil, err
	}

	return &pb.SendMessageResponse{
		Message: "Message has been sent",
	}, nil
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		if err := c.Conn.WriteJSON(message); err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.Hub.Unregister <- c.Participant
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

		msg := &models.Message{
			ParticipantID: c.Participant.ID,
			Text:          string(m),
		}

		c.Hub.Broadcast <- msg
	}
}

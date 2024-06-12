package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/temuka-messaging-service/config"
	"github.com/temuka-messaging-service/models"
	"github.com/temuka-messaging-service/pb"
)

func (s *server) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	db := config.GetDBInstance()

	newConversation := models.Conversation{
		Title:  req.Title,
		UserID: int(req.CreatorId),
	}

	if err := db.Create(&newConversation).Error; err != nil {
		return nil, err
	}

	return &pb.CreateConversationResponse{
		Message: "Conversation has been created",
		Data: &pb.Conversation{
			Id:     int32(newConversation.ID),
			Title:  newConversation.Title,
			UserId: int32(newConversation.UserID),
		},
	}, nil
}

func (s *server) JoinConversation(ctx context.Context, req *pb.JoinConversationRequest) (*pb.JoinConversationResponse, error) {
	db := config.GetDBInstance()

	newParticipant := models.Participant{
		ConversationID: int(req.Id),
		UserID:         int(req.UserId),
	}

	if err := db.Create(&newParticipant).Error; err != nil {
		return nil, err
	}

	client := &Client{
		Message:     make(chan *models.Message, 10),
		Participant: &newParticipant,
		Hub:         hub,
		DB:          db,
	}

	hub.Register <- &newParticipant

	go client.writeMessage()
	client.readMessage()

	return &pb.JoinConversationResponse{
		Message: "Joined conversation successfully",
	}, nil
}

func (s *server) GetConversations(ctx context.Context, req *pb.GetConversationsRequest) (*pb.GetConversationsResponse, error) {
	db := config.GetDBInstance()

	var conversations []models.Conversation
	err := db.Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("participants.id = ?", req.ParticipantId).
		Preload("Participants").
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}

	var pbConversations []*pb.Conversation
	for _, conv := range conversations {
		pbConversations = append(pbConversations, &pb.Conversation{
			Id:     int32(conv.ID),
			Title:  conv.Title,
			UserId: int32(conv.UserID),
		})
	}

	return &pb.GetConversationsResponse{
		Message: "Conversations retrieved successfully",
		Data:    pbConversations,
	}, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

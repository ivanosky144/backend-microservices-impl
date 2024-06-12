package handlers

import (
	"github.com/temuka-messaging-service/pb"
	"github.com/temuka-messaging-service/ws"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
}

// NewServer creates a new gRPC server.
func NewServer() *server {
	return &server{}
}

var hub *ws.Hub
var db *gorm.DB

package handlers

import (
	"github.com/temuka-content-service/pb"
)

type server struct {
	pb.UnimplementedContentServiceServer
}

// NewServer creates a new gRPC server.
func NewServer() *server {
	return &server{}
}

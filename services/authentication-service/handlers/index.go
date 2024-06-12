package handlers

import (
	"github.com/temuka-authentication-service/pb"
)

type server struct {
	pb.UnimplementedAuthenticationServiceServer
}

// NewServer creates a new gRPC server.
func NewServer() *server {
	return &server{}
}

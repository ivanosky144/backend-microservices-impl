package main

import (
	"log"
	"net"

	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/handlers"
	"github.com/temuka-content-service/models"
	"github.com/temuka-content-service/pb"
	"google.golang.org/grpc"
)

func main() {
	config.OpenConnection()
	var db = config.GetDBInstance()

	if db == nil {
		log.Fatal("Database connection is nil")
	}

	if err := db.AutoMigrate(&models.Community{}, &models.CommunityMember{}, &models.Comment{}, &models.Post{}, &models.UserLike{}, &models.UserVote{}, &models.Moderator{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Auto-migration completed.")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC service implementation
	pb.RegisterContentServiceServer(grpcServer, handlers.NewServer())

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is listening on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

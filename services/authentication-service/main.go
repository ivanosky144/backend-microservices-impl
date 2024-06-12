package main

import (
	"log"
	"net"

	"github.com/temuka-authentication-service/config"
	"github.com/temuka-authentication-service/handlers"
	"github.com/temuka-authentication-service/models"
	"github.com/temuka-authentication-service/pb"
	"google.golang.org/grpc"
)

func main() {
	config.OpenConnection()
	var db = config.GetDBInstance()

	if db == nil {
		log.Fatal("Database connection is nil")
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Auto-migration completed.")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC service implementation
	pb.RegisterAuthenticationServiceServer(grpcServer, handlers.NewServer())

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

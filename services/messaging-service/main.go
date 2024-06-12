package main

import (
	"log"
	"net"

	"github.com/temuka-messaging-service/config"
	"github.com/temuka-messaging-service/handlers"
	"github.com/temuka-messaging-service/models"
	"github.com/temuka-messaging-service/pb"
	"github.com/temuka-messaging-service/ws"
	"google.golang.org/grpc"
)

func main() {
	config.OpenConnection()
	var db = config.GetDBInstance()

	if db == nil {
		log.Fatal("Database connection is nil")
	}

	if err := db.AutoMigrate(&models.Conversation{}, &models.Message{}, &models.Participant{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Auto-migration completed.")

	hub := ws.NewHub()
	go hub.Run()

	ws.Init(hub, db)

	grpcServer := grpc.NewServer()

	pb.RegisterMessagingServiceServer(grpcServer, handlers.NewServer())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

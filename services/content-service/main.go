package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/models"
	"github.com/temuka-content-service/routes"

	"gorm.io/gorm"
)

func EnableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
	config.OpenConnection()
	var db *gorm.DB = config.GetDBInstance()

	if config.Database == nil {
		log.Fatal("Database connection is nil")
	}
	if err := config.Database.AutoMigrate(&models.Community{}, &models.Post{}, &models.Comment{}, &models.CommunityMember{}, &models.CommunityPost{}, &models.Moderator{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Printf("Database : %v", db)
	log.Println("Auto-migration completed.")
	router := mux.NewRouter()

	router.PathPrefix("/api/post").Handler(http.StripPrefix("/api/post", routes.PostRoutes()))
	router.PathPrefix("/api/community").Handler(http.StripPrefix("/api/community", routes.CommunityRoutes()))
	router.PathPrefix("/api/comment").Handler(http.StripPrefix("/api/comment", routes.CommentRoutes()))

	http.Handle("/", router)

	log.Println("Content service is listening on port 3400")
	log.Fatal(http.ListenAndServe("localhost:3400", nil))
}

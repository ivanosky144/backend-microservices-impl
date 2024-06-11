package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temuka-authentication-service/config"
	"github.com/temuka-authentication-service/models"
	"github.com/temuka-authentication-service/routes"

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
	if err := config.Database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Printf("Database : %v", db)
	log.Println("Auto-migration completed.")
	router := mux.NewRouter()

	router.PathPrefix("/api/auth").Handler(http.StripPrefix("/api/auth", routes.AuthRoutes()))
	router.PathPrefix("/api/user").Handler(http.StripPrefix("/api/user", routes.UserRoutes()))

	http.Handle("/", router)

	log.Println("Authentication service is listening on port 3300")
	log.Fatal(http.ListenAndServe("localhost:3300", nil))
}

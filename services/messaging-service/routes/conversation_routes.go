package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-messaging-service/handlers"
)

func ConversationRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/conversation", handlers.CreateConversation).Methods("POST")

	return r
}

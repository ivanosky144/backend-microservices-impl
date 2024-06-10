package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-content-service/handlers"
)

func CommunityRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.CreateCommunity).Methods("POST")
	r.HandleFunc("/", handlers.JoinCommunity).Methods("POST")

	return r
}

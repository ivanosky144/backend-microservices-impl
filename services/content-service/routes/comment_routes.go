package routes

import (
	"github.com/gorilla/mux"
	"github.com/temuka-content-service/handlers"
)

func CommentRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.AddComment).Methods("POST")

	return r
}

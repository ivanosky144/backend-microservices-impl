package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temuka-messaging-service/routes"
)

func EnableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/auth").Handler(routes.ConversationRoutes())

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8003", nil))
}

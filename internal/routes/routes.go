package routes

import (
	"imantz/daily_it_meeting_helper_ppp/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ws", handlers.WsEndpoint)
	router.HandleFunc("/generate", handlers.GenerateHandler).Methods("POST")
	router.HandleFunc("/current-entry", handlers.CurrentEntryHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}

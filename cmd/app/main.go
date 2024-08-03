package main

import (
	"imantz/daily_it_meeting_helper_ppp/internal/routes"
	"imantz/daily_it_meeting_helper_ppp/internal/services"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = services.LoadEntries()
	if err != nil {
		log.Fatalf("Error loading entries: %v", err)
	}

	router := routes.SetupRouter()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

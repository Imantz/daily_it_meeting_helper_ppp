package handlers

import (
	"encoding/json"
	"fmt"
	"imantz/daily_it_meeting_helper_ppp/internal/models"
	"imantz/daily_it_meeting_helper_ppp/internal/services"
	"net/http"
	"time"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	formattedText := fmt.Sprintf("Progress: %s\nPlans: %s\nProblems: %s",
		msg.Progress, msg.Plans, msg.Problems)

	response := map[string]string{
		"formattedText": formattedText,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CurrentEntryHandler(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().Format("2006-01-02")
	entry, exists := services.EntriesByDate[currentDate]
	if !exists {
		entry = models.Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

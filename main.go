package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Progress string `json:"progress"`
	Plans    string `json:"plans"`
	Problems string `json:"problems"`
}

var entriesByDate = make(map[string]Message)

const dataFile = "entries.json"

func saveEntries() error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(entriesByDate)
}

func loadEntries() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&entriesByDate)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		currentDate := time.Now().Format("2006-01-02")
		entriesByDate[currentDate] = msg

		err = saveEntries()
		if err != nil {
			log.Println("Error saving entries:", err)
		}

		log.Printf("Saved entry for date %s: %+v\n", currentDate, msg)
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Simulate calling ChatGPT API
	formattedText := fmt.Sprintf("Progress: %s\nPlans: %s\nProblems: %s",
		msg.Progress, msg.Plans, msg.Problems)

	response := map[string]string{
		"formattedText": formattedText,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func currentEntryHandler(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().Format("2006-01-02")
	entry, exists := entriesByDate[currentDate]
	if !exists {
		entry = Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func main() {
	err := loadEntries()
	if err != nil {
		log.Fatalf("Error loading entries: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/ws", wsEndpoint)
	router.HandleFunc("/generate", generateHandler).Methods("POST")
	router.HandleFunc("/current-entry", currentEntryHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", router)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

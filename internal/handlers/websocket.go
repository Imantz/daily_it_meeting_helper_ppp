package handlers

import (
	"imantz/daily_it_meeting_helper_ppp/internal/models"
	"imantz/daily_it_meeting_helper_ppp/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		currentDate := time.Now().Format("2006-01-02")
		services.EntriesByDate[currentDate] = msg

		err = services.SaveEntries()
		if err != nil {
			log.Println("Error saving entries:", err)
		}

		log.Printf("Saved entry for date %s: %+v\n", currentDate, msg)
	}
}

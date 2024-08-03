package services

import (
	"encoding/json"
	"imantz/daily_it_meeting_helper_ppp/internal/models"
	"os"
)

const dataFile = "entries.json"

var EntriesByDate = make(map[string]models.Message)

func SaveEntries() error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(EntriesByDate)
}

func LoadEntries() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&EntriesByDate)
}

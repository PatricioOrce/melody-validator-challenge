package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"melody-validator-challenge/application"
	"melody-validator-challenge/cmd/server/models"
	"net/http"
)

func ValidateMelodyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var body models.MelodyRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.Melody == "" {
		http.Error(w, "Melody field is required", http.StatusBadRequest)
		return
	}

	errIndex := application.ValidateMelody(body.Melody)
	if errIndex > -1 {
		errResponse := models.ErrorResponse{
			Cause: fmt.Sprintf("error at position %d", errIndex),
		}
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(errResponse); encodeErr != nil {
			log.Printf("Error encoding error response: %v", encodeErr)
		}
		return
	}
	response := application.MapMelody(body.Melody)
	if encodeErr := json.NewEncoder(w).Encode(&response); encodeErr != nil {
		log.Printf("Error encoding error response: %v", encodeErr)
	}
	w.WriteHeader(http.StatusOK)
}

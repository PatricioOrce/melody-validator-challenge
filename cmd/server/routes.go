package server

import (
	"melody-validator-challenge/cmd/server/handlers"
	"net/http"
)

func initRoutes() {
	http.HandleFunc("/melody/validate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		handlers.ValidateMelodyHandler(w, r)
	})
}

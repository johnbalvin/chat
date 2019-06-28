package main

import (
	"chat/backend/userlogs"
	"net/http"

	"github.com/gorilla/mux"
)

func handleMisStadistics(r *mux.Router) {
	r.HandleFunc("/chgCo", handleDarkBlockers).Methods("POST")
}

//handleDarkBlockers it will handle dark mode and blockers
func handleDarkBlockers(w http.ResponseWriter, r *http.Request) {
	userlogs.SaveBloquer(r)
}

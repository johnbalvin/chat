package main

import (
	"chat/backend/chat"
	"chat/backend/errors"
	"chat/backend/sessions"
	"chat/backend/users"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type smallInfo struct {
	UserID string `json:"i"`
	Photo  string `json:"p"`
	Name   string `json:"n"`
}

func handleChat(r *mux.Router) {
	r.HandleFunc("/", chatHandler)
	r.HandleFunc("/UsersData", usersData).Methods("POST")
}

func usersData(w http.ResponseWriter, r *http.Request) {
	usersIdsRaw := r.FormValue("u")
	usersIDs := make(map[string]bool)
	if err := json.Unmarshal([]byte(usersIdsRaw), &usersIDs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("main -> usersData:1 -> err: ", err)
		return
	}
	users, err := users.ByIDs(usersIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("main -> usersData:2 -> err: ", err)
		return
	}
	var infoToSend []smallInfo
	for _, user := range users { //removin inncesary information
		infoToSend = append(infoToSend, smallInfo{UserID: user.ID, Photo: user.Photo, Name: user.Name})
	}
	if err := json.NewEncoder(w).Encode(infoToSend); err != nil {
		log.Println("main -> usersData:3 -> err:", err)
		return
	}
}
func chatHandler(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	var user users.User
	if userID != "" {
		var err error
		user, err = sessions.Confirm(userID, sessionID)
		if err != nil {
			sessions.SetCookie(trackerID, "", "", w)
			log.Println("main -> chatHandler:1 -> err:", err)
			return
		}
	}
	if r.Method == "GET" {
		data := struct {
			User     users.User
			Messages []chat.Message
		}{User: user}
		msgs, _ := chat.GetMsgs(20)
		data.Messages = msgs
		if err := servirHTMLS.ExecuteTemplate(w, "frontEnd/chat/index.html", data); err != nil {
			log.Println("main -> chatHandler:2 -> err:", err)
			return
		}
	} else if r.Method == "POST" {
		c := r.FormValue("c")
		switch c {
		case "n": //new message
			if userID == "" { //means no user
				return
			}
			m := r.FormValue("m")
			_, err := chat.AddMsg(user, m)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusPreconditionFailed)
				log.Println("main -> chatHandler:3 -> err:", err)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> chatHandler:4 -> err:", err)
				return
			}
		case "b": //search message before
			unixNanoString := r.FormValue("w")
			unixNano, err := strconv.ParseInt(unixNanoString, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> chatHandler:5 -> err:", err)
				return
			}
			messages, err := chat.GetMsgsBeforeTime(unixNano, 20)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> chatHandler:6 -> err:", err)
				return
			}
			if len(messages) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			data := struct {
				ContinueSearching bool           `json:"s"`
				Messages          []chat.Message `json:"m"`
			}{Messages: messages}
			if len(messages) == 20 {
				data.ContinueSearching = true
			}
			if err := json.NewEncoder(w).Encode(data); err != nil {
				log.Println("main -> chatHandler:7 -> err:", err)
				return
			}
		}
	}
}

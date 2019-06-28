package main

import (
	"chat/backend/sessions"
	"chat/backend/users"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func handleSettings(r *mux.Router) {
	r.HandleFunc("/Settings", settings)
}
func settings(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2e+6)
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	if userID == "" {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	user, err := sessions.Confirm(userID, sessionID)
	if err != nil {
		sessions.SetCookie(trackerID, "", "", w)
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		log.Println("main -> settings:1 -> err:", err)
		return
	}
	if r.Method == "GET" {
		servir := struct {
			User users.User
		}{User: user}
		if err := servirHTMLS.ExecuteTemplate(w, "frontEnd/settings/index.html", servir); err != nil {
			log.Println("main -> settings:2 -> err: ", err)
			return
		}
	} else if r.Method == "POST" {
		name := html.EscapeString(strings.TrimSpace(r.FormValue("n")))
		if len(name) != 0 {
			if len(name) >= 70 {
				log.Println("main -> settings:3 -> name: ", name, " userID: ", user.ID)
				return
			}
			if err := user.ChangeName(name); err != nil {
				log.Println("main -> settings:4 -> err: ", err)
				return
			}
		}
		file, header, err := r.FormFile("i")
		if err == http.ErrMissingFile {
			return
		}
		if err != nil {
			log.Println("main -> settings:5 -> err: ", err, " userID: ", user.ID)
			return
		}
		if _, err := user.ChangesPhoto(file, header); err != nil {
			log.Println("main -> settings:6 -> err: ", err, " userID: ", user.ID)
			return
		}
	}
}

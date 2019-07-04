package main

import (
	"chat/backend/errors"
	"chat/backend/sessions"
	"chat/backend/users"
	"chat/backend/users/recover"
	"chat/backend/users/register"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func handleAcces(r *mux.Router) {
	r.HandleFunc("/login", login)
	r.HandleFunc("/close", closeSession).Methods("GET")
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/reset", resetpassword)
}
func login(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	if r.Method == "GET" {
		if userID == "" {
			w.Write(sevirHTMLSinData["frontEnd/acces/login/index.html"])
			return
		}
		_, err := sessions.Confirm(userID, sessionID)
		if err != nil {
			sessions.SetCookie(trackerID, "", "", w)
			log.Println("main -> login:1 -> err:", err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "POST" {
		email := html.EscapeString(strings.ToLower(strings.TrimSpace(r.FormValue("email"))))
		if email == "" {
			return
		}
		pass := r.FormValue("pass")
		user, err := users.ByEmail(email)
		if err == errors.Empty {
			w.WriteHeader(http.StatusNotFound)
			log.Println("main -> login:2 -> email:", email)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("main -> login:3 -> err:", err)
			return
		}
		isEqual, err := user.PassCompare(pass)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("main -> login:4 -> err:", err)
			return
		}
		if !isEqual {
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Println("main -> login:5 -> not equal")
			return
		}
		if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("main -> login:6 -> err:", err)
			return
		}
	}
}
func signup(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	if r.Method == "GET" {
		if userID == "" {
			w.Write(sevirHTMLSinData["frontEnd/acces/signup/index.html"])
			return
		}
		_, err := sessions.Confirm(userID, sessionID)
		if err != nil {
			sessions.SetCookie(trackerID, "", "", w)
			log.Println("main -> signup:1 -> err:", err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "POST" {
		email := html.EscapeString(strings.ToLower(strings.TrimSpace(r.FormValue("email"))))
		if email == "" {
			return
		}
		step := r.FormValue("s")
		switch step {
		case "1":
			err := register.Step1(email)
			if err == errors.Duplicated {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> signup:2 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> signup:3 -> err:", err)
				return
			}
		case "2":
			codeEmail := r.FormValue("code")
			err := register.Step2(email, codeEmail)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> signup:4 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> signup:5 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> signup:6 -> err:", err)
				return
			}
		case "3":
			name := html.EscapeString(strings.TrimSpace(r.FormValue("name")))
			pass := r.FormValue("pass")
			code := r.FormValue("code")
			if len(name) < 0 || len(pass) > 100 { //se supone que esto ya se lo valida desde el cliente
				log.Println("main -> signup:7 -> name:", name, " pass: ", pass)
				return
			}
			user, err := register.Step3(email, code, pass, name)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> signup:8 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> signup:9 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> signup:10 -> err:", err)
				return
			}
			if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> signup:11 -> err:", err)
				return
			}
		}
	}
}
func resetpassword(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, _, _ := sessions.ReadAndAsignCookie(w, r)
	if r.Method == "GET" {
		if userID == "" {
			w.Write(sevirHTMLSinData["frontEnd/acces/reset/index.html"])
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "POST" {
		email := html.EscapeString(strings.ToLower(strings.TrimSpace(r.FormValue("email"))))
		step := r.FormValue("s")
		switch step {
		case "1":
			err := recover.Step1(email)
			if err == errors.Empty {
				w.WriteHeader(http.StatusNoContent)
				log.Println("main -> resetpassword:1 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> resetpassword:2 -> err:", err)
				return
			}
		case "2":
			codeEmail := r.FormValue("code")
			err := recover.Step2(email, codeEmail)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> resetpassword:3 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> resetpassword:4 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> resetpassword:5 -> err:", err)
				return
			}
		case "3":
			code := r.FormValue("code")
			contra := r.FormValue("pass")
			user, err := recover.Step3(email, code, contra)
			if err == errors.Empty {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> resetpassword:6 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> resetpassword:7 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> resetpassword:8 -> err:", err)
				return
			}
			if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> resetpassword:9 -> err:", err)
				return
			}
		}
	}
}
func closeSession(w http.ResponseWriter, r *http.Request) {
	sessions.Close(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

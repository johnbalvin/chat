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
	r.HandleFunc("/", index)
	r.HandleFunc("/Login", login)
	r.HandleFunc("/Close", closeSession).Methods("GET")
	r.HandleFunc("/Register", registerAccount)
	r.HandleFunc("/ResetPassword", recoverAccount)
}
func index(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	if userID == "" {
		w.Write(sevirHTMLSinData["frontEnd/acces/index/index.html"])
		return
	}
	_, err := sessions.Confirm(userID, sessionID)
	if err != nil {
		sessions.SetCookie(trackerID, "", "", w)
		log.Println("main -> index:1 -> err:", err)
		return
	}
	http.Redirect(w, r, "/Chat", http.StatusSeeOther)
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
		http.Redirect(w, r, "/Chat", http.StatusSeeOther)
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
			log.Println("main -> login:5 -> no es igual")
			return
		}
		if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("main -> login:6 -> err:", err)
			return
		}
	}
}
func registerAccount(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, sessionID, _ := sessions.ReadAndAsignCookie(w, r)
	if r.Method == "GET" {
		if userID == "" {
			w.Write(sevirHTMLSinData["frontEnd/acces/register/index.html"])
			return
		}
		_, err := sessions.Confirm(userID, sessionID)
		if err != nil {
			sessions.SetCookie(trackerID, "", "", w)
			log.Println("main -> registerAccount:1 -> err:", err)
			return
		}
		http.Redirect(w, r, "/Chat", http.StatusSeeOther)
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
				log.Println("main -> registerAccount:2 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> registerAccount:3 -> err:", err)
				return
			}
		case "2":
			codeEmail := r.FormValue("code")
			err := register.Step2(email, codeEmail)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> registerAccount:4 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> registerAccount:5 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> registerAccount:6 -> err:", err)
				return
			}
		case "3":
			name := html.EscapeString(strings.TrimSpace(r.FormValue("name")))
			pass := r.FormValue("pass")
			code := r.FormValue("code")
			if len(name) < 0 || len(pass) > 100 { //se supone que esto ya se lo valida desde el cliente
				log.Println("main -> registerAccount:7 -> name:", name, " pass: ", pass)
				return
			}
			user, err := register.Step3(email, code, pass, name)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> registerAccount:8 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> registerAccount:9 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> registerAccount:10 -> err:", err)
				return
			}
			if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> registerAccount:11 -> err:", err)
				return
			}
		}
	}
}
func recoverAccount(w http.ResponseWriter, r *http.Request) {
	trackerID, userID, _, _ := sessions.ReadAndAsignCookie(w, r)
	if r.Method == "GET" {
		if userID == "" {
			w.Write(sevirHTMLSinData["frontEnd/acces/recover/index.html"])
			return
		}
		http.Redirect(w, r, "/Chat", http.StatusSeeOther)
	} else if r.Method == "POST" {
		correo := html.EscapeString(strings.ToLower(strings.TrimSpace(r.FormValue("correo"))))
		paso := r.FormValue("p")
		switch paso {
		case "1":
			err := recover.Step1(correo)
			if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> recoverAccount:1 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> recoverAccount:2 -> err:", err)
				return
			}
		case "2":
			codeEmail := r.FormValue("code")
			err := recover.Step2(correo, codeEmail)
			if err == errors.TimesOut {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> recoverAccount:3 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> recoverAccount:4 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> recoverAccount:5 -> err:", err)
				return
			}
		case "3":
			code := r.FormValue("code")
			contra := r.FormValue("pass")
			user, err := recover.Step3(correo, code, contra)
			if err == errors.Empty {
				w.WriteHeader(http.StatusForbidden)
				log.Println("main -> recoverAccount:6 -> err:", err)
				return
			} else if err == errors.Empty {
				w.WriteHeader(http.StatusConflict)
				log.Println("main -> recoverAccount:7 -> err:", err)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> recoverAccount:8 -> err:", err)
				return
			}
			if _, _, err := sessions.Start(trackerID, user, w, r); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("main -> recoverAccount:9 -> err:", err)
				return
			}
		}
	}
}
func closeSession(w http.ResponseWriter, r *http.Request) {
	sessions.Close(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

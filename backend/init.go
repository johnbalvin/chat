package main

import (
	"chat/backend/servetousers"
	"chat/backend/userlogs"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var servirHTMLS *template.Template
var servirSW = make(map[string][]byte)
var sevirHTMLSinData = make(map[string][]byte)

func init() {
	servirHTMLS, sevirHTMLSinData, _ = servetousers.Get()
}

func app() {
	me := mux.NewRouter()
	registerHandlers(me)
	srv := &http.Server{
		Addr: ":8070",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			me.ServeHTTP(w, r)
			go userlogs.SaveNormal(r)
		}),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}
	fmt.Println("---Server started at port: 8070--")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err) // Error starting or closing listener:
	}
}
func main() {
	app()
}

func registerHandlers(r *mux.Router) {
	handleAcces(r)
	handleSettings(r)
	handleChat(r)
	handleMisStadistics(r)
}
